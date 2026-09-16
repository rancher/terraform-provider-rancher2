#!/usr/bin/env node

/**
 * Sensitive Documentation Warning Updater
 * 
 * Scans Go provider schemas to identify all fields marked as `Sensitive: true`
 * and programmatically appends a security warning to their descriptions
 * inside the markdown documentation files (resources and data-sources).
 */

import fs from 'fs/promises';
import path from 'path';

// Default warning text to append to sensitive fields
const WARNING_TEXT = ' - **Note:** This value is stored in cleartext in Terraform state regardless of the `sensitive` flag. It is highly recommended to use a state backend with encryption and restricted access.';

function printHelp() {
    console.log(`
Usage: node update_sensitive_docs.js [options]

Options:
  -h, --help        Show this help menu.
  -d, --dry-run     Run the scanner without writing changes to the markdown files.
  --go-dir <path>   Path to the Go schemas directory (default: ./rancher2).
  --docs-dir <path> Path to the docs directory containing resources and data-sources (default: ./docs).
`);
}

/**
 * Parses Go files to extract all sensitive field names.
 * Matches fields like "secret_key": { ... Sensitive: true ... }
 */
async function scanGoSchemas(goDir) {
    const sensitiveFields = new Set();
    let files;
    try {
        files = await fs.readdir(goDir);
    } catch (err) {
        if (err.code === 'ENOENT') {
            return sensitiveFields;
        }
        throw err;
    }

    // Regex to find field blocks in Go maps: "field_name": {
    const fieldRegex = /"([a-zA-Z0-9_]+)"\s*:\s*\{/g;

    for (const file of files) {
        if (path.extname(file) !== '.go') continue;

        const filePath = path.join(goDir, file);
        const content = await fs.readFile(filePath, 'utf-8');

        let match;
        // Reset regex index
        fieldRegex.lastIndex = 0;

        while ((match = fieldRegex.exec(content)) !== null) {
            const fieldName = match[1];
            const startIndex = match.index;

            // Find the end of this field's block by counting braces
            let braceCount = 0;
            let foundStart = false;
            let blockContent = '';

            for (let i = startIndex; i < content.length; i++) {
                const char = content[i];
                if (char === '{') {
                    braceCount++;
                    foundStart = true;
                } else if (char === '}') {
                    braceCount--;
                }

                if (foundStart) {
                    blockContent += char;
                }

                if (foundStart && braceCount === 0) {
                    break;
                }
            }

            // If the block contains Sensitive: true, mark the field as sensitive
            if (/\bSensitive\s*:\s*true\b/.test(blockContent)) {
                sensitiveFields.add(fieldName);
            }
        }
    }

    return sensitiveFields;
}

/**
 * Programmatically updates markdown files in resources and data-sources.
 */
async function updateMarkdownDocs(docsDir, sensitiveFields, dryRun) {
    let updateCount = 0;

    const subDirs = ['resources', 'data-sources'];
    for (const subDir of subDirs) {
        const subDirPath = path.join(docsDir, subDir);
        let files;
        try {
            files = await fs.readdir(subDirPath);
        } catch (err) {
            if (err.code === 'ENOENT') continue;
            throw err;
        }

        for (const file of files) {
            if (path.extname(file) !== '.md') continue;

            const filePath = path.join(subDirPath, file);
            const content = await fs.readFile(filePath, 'utf-8');
            const lines = content.split('\n');
            let modified = false;

            const updatedLines = lines.map(line => {
                // Match bullet points starting with: * `field_name`, - `field_name`, etc.
                const match = line.match(/^\s*[\*-]\s*[`*]([a-zA-Z0-9_]+)[`*]/);
                if (match) {
                    const fieldName = match[1];
                    if (sensitiveFields.has(fieldName)) {
                        // Check if warning already exists to prevent duplicate appends
                        if (!line.includes('stored in cleartext')) {
                            modified = true;
                            updateCount++;
                            return line + WARNING_TEXT;
                        }
                    }
                }
                return line;
            });

            if (modified && !dryRun) {
                await fs.writeFile(filePath, updatedLines.join('\n'), 'utf-8');
                console.log(`Updated sensitive field documentation in: ${filePath}`);
            } else if (modified && dryRun) {
                console.log(`[DRY RUN] Would update sensitive field documentation in: ${filePath}`);
            }
        }
    }

    return updateCount;
}

async function main() {
    const args = process.argv.slice(2);
    let dryRun = false;
    let goDir = './rancher2';
    let docsDir = './docs';

    for (let i = 0; i < args.length; i++) {
        const arg = args[i];
        if (arg === '-h' || arg === '--help') {
            printHelp();
            process.exit(0);
        } else if (arg === '-d' || arg === '--dry-run') {
            dryRun = true;
        } else if (arg === '--go-dir') {
            if (i + 1 >= args.length) {
                console.error('Error: --go-dir requires a path argument');
                process.exit(1);
            }
            goDir = args[++i];
        } else if (arg === '--docs-dir') {
            if (i + 1 >= args.length) {
                console.error('Error: --docs-dir requires a path argument');
                process.exit(1);
            }
            docsDir = args[++i];
        } else {
            console.error(`Unknown option: ${arg}`);
            printHelp();
            process.exit(1);
        }
    }

    console.log('Scanning Go schemas for sensitive fields...');
    const sensitiveFields = await scanGoSchemas(goDir);
    console.log(`Found ${sensitiveFields.size} unique sensitive fields:`, Array.from(sensitiveFields).join(', '));

    console.log('\nScanning and updating documentation files...');
    const updateCount = await updateMarkdownDocs(docsDir, sensitiveFields, dryRun);

    if (dryRun) {
        console.log(`\n[DRY RUN] Completed. Would update ${updateCount} field descriptions.`);
    } else {
        console.log(`\nCompleted. Updated ${updateCount} field descriptions across documentation files.`);
    }

    // Exit with code 1 if there are findings/updates, and code 0 otherwise
    if (updateCount > 0) {
        process.exit(1);
    } else {
        process.exit(0);
    }
}

main().catch(err => {
    console.error('Fatal error during execution:', err);
    process.exit(1);
});
