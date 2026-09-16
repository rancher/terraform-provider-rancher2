package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Pre-compiled regular expressions for performance.
var (
	blockHeaderRegex  = regexp.MustCompile(`^\s*(variable|output)\s+"([^"]+)"\s*\{`)
	sensitiveValRegex = regexp.MustCompile(`\bsensitive\s*=\s*true\b`)
)

// IgnoreList holds global and file-specific ignore mappings.
type IgnoreList struct {
	global map[string]bool
	local  map[string]map[string]bool
}

// Check if a field name matches any target keywords.
func matchesKeywords(name string, keywords []string) bool {
	lower := strings.ToLower(name)
	for _, kw := range keywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

// loadKeywords reads and parses the keywords configuration file.
func loadKeywords(path string) ([]string, error) {
	defaultKeywords := []string{"password", "token", "credential", "pemencodedprivatekey", "secret", "private"}
	if path == "" {
		return defaultKeywords, nil
	}
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return defaultKeywords, nil // Fall back to default keywords if file doesn't exist.
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read keywords file %s: %w", path, err)
	}
	lines := strings.Split(string(content), "\n")
	var keywords []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		keywords = append(keywords, strings.ToLower(line))
	}
	if len(keywords) == 0 {
		return defaultKeywords, nil
	}
	return keywords, nil
}

// loadIgnoreList reads and parses the ignore file.
func loadIgnoreList(path string) (*IgnoreList, error) {
	ignore := &IgnoreList{
		global: make(map[string]bool),
		local:  make(map[string]map[string]bool),
	}
	if path == "" {
		return ignore, nil
	}
	content, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return ignore, nil // default to empty ignore list if file doesn't exist.
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read ignore file %s: %w", path, err)
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.Contains(line, ":") {
			parts := strings.SplitN(line, ":", 2)
			file := strings.TrimSpace(parts[0])
			field := strings.TrimSpace(parts[1])
			if ignore.local[file] == nil {
				ignore.local[file] = make(map[string]bool)
			}
			ignore.local[file][field] = true
		} else {
			ignore.global[line] = true
		}
	}
	return ignore, nil
}

// IsIgnored returns true if a field is ignored globally or for the specific file.
func (i *IgnoreList) IsIgnored(filePath, fieldName string) bool {
	if i.global[fieldName] {
		return true
	}
	if fields, ok := i.local[filePath]; ok {
		if fields[fieldName] {
			return true
		}
	}
	return false
}

// Find schema maps in Go files and analyze them.
func checkGoFile(fset *token.FileSet, file *ast.File, filePath string, ignoreList *IgnoreList, keywords []string) []string {
	var findings []string

	ast.Inspect(file, func(n ast.Node) bool {
		compLit, ok := n.(*ast.CompositeLit)
		if !ok {
			return true
		}

		// Check if it's a map type, e.g. map[string]*schema.Schema or map[string]schema.Schema
		mapType, ok := compLit.Type.(*ast.MapType)
		if !ok {
			return true
		}

		// Ensure key is string
		keyIdent, ok := mapType.Key.(*ast.Ident)
		if !ok || keyIdent.Name != "string" {
			return true
		}

		// We check if value is schema.Schema or *schema.Schema
		isSchemaMap := false
		switch valType := mapType.Value.(type) {
		case *ast.StarExpr: // *schema.Schema
			if isSchemaType(valType.X) {
				isSchemaMap = true
			}
		case *ast.SelectorExpr: // schema.Schema
			if isSchemaType(valType) {
				isSchemaMap = true
			}
		}

		if !isSchemaMap {
			return true
		}

		// Traverse map elements
		for _, elt := range compLit.Elts {
			kv, ok := elt.(*ast.KeyValueExpr)
			if !ok {
				continue
			}

			// Extract map key (field name)
			lit, ok := kv.Key.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				continue
			}
			fieldName := strings.Trim(lit.Value, "\"`")

			if !matchesKeywords(fieldName, keywords) {
				continue
			}

			if ignoreList.IsIgnored(filePath, fieldName) {
				continue
			}

			// Now inspect the value (struct literal schema.Schema)
			structLit, ok := getCompositeLit(kv.Value)
			if !ok {
				continue
			}

			// Check if Sensitive: true is defined
			hasSensitive := false
			for _, structElt := range structLit.Elts {
				structKV, ok := structElt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}

				ident, ok := structKV.Key.(*ast.Ident)
				if !ok {
					continue
				}

				if ident.Name == "Sensitive" {
					if valIdent, ok := structKV.Value.(*ast.Ident); ok && valIdent.Name == "true" {
						hasSensitive = true
					}
				}
			}

			if !hasSensitive {
				pos := fset.Position(kv.Key.Pos())
				findings = append(findings, fmt.Sprintf("%s:%d: Go schema attribute %q is not marked as sensitive", filePath, pos.Line, fieldName))
			}
		}

		return true
	})

	return findings
}

func isSchemaType(expr ast.Expr) bool {
	selector, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := selector.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "schema" && selector.Sel.Name == "Schema"
}

func getCompositeLit(expr ast.Expr) (*ast.CompositeLit, bool) {
	switch e := expr.(type) {
	case *ast.CompositeLit:
		return e, true
	case *ast.UnaryExpr: // check for &schema.Schema{...}
		if e.Op == token.AND {
			return getCompositeLit(e.X)
		}
	}
	return nil, false
}

// Check Terraform (.tf) files for unmarked sensitive variable and output blocks.
func checkTfFile(filePath string, ignoreList *IgnoreList, keywords []string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TF file %s: %w", filePath, err)
	}
	return analyzeTfContent(string(content), filePath, ignoreList, keywords)
}

// Helper to count braces while ignoring characters inside double quotes.
func countBraces(s string) (open, closedBraces int) {
	inQuote := false
	escaped := false
	for _, r := range s {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if !inQuote {
			switch r {
			case '{':
				open++
			case '}':
				closedBraces++
			}
		}
	}
	return open, closedBraces
}

// Helper to strip comments safely while respecting double quotes.
func stripComments(line string) string {
	inQuote := false
	escaped := false
	for i, r := range line {
		if escaped {
			escaped = false
			continue
		}
		if r == '\\' {
			escaped = true
			continue
		}
		if r == '"' {
			inQuote = !inQuote
			continue
		}
		if !inQuote {
			if r == '#' {
				return line[:i]
			}
			if i < len(line)-1 && line[i:i+2] == "//" {
				// Ensure it's not part of a URL (e.g., http:// or https://)
				if i == 0 || line[i-1] != ':' {
					return line[:i]
				}
			}
		}
	}
	return line
}

func analyzeTfContent(content string, filePath string, ignoreList *IgnoreList, keywords []string) ([]string, error) {
	lines := strings.Split(content, "\n")
	var findings []string

	type Block struct {
		blockType string
		name      string
		startLine int
		braces    int
		sensitive bool
	}

	var currentBlock *Block
	inHeredoc := false
	heredocDelim := ""
	inBlockComment := false

	for i, line := range lines {
		lineNum := i + 1
		trimmed := strings.TrimSpace(line)

		// Handle block comments (/* ... */)
		if inBlockComment {
			if idx := strings.Index(line, "*/"); idx != -1 {
				inBlockComment = false
				line = line[idx+2:]
				trimmed = strings.TrimSpace(line)
			} else {
				continue
			}
		}

		// Check for block comment start
		if idx := strings.Index(line, "/*"); idx != -1 {
			// Check if it's closed on the same line
			if closeIdx := strings.Index(line[idx+2:], "*/"); closeIdx != -1 {
				line = line[:idx] + line[idx+2+closeIdx+2:]
				trimmed = strings.TrimSpace(line)
			} else {
				inBlockComment = true
				line = line[:idx]
				trimmed = strings.TrimSpace(line)
			}
		}

		// Handle heredoc state
		if inHeredoc {
			if strings.HasSuffix(trimmed, heredocDelim) || trimmed == heredocDelim {
				inHeredoc = false
				heredocDelim = ""
			}
			continue
		}

		// Strip comments first to get activePart for syntax analysis
		activePart := stripComments(line)

		if currentBlock == nil {
			// Look for new block
			matches := blockHeaderRegex.FindStringSubmatch(line)
			if len(matches) == 3 {
				currentBlock = &Block{
					blockType: matches[1],
					name:      matches[2],
					startLine: lineNum,
					braces:    1,
					sensitive: false,
				}
				// Also count braces in the remainder of the line after the header
				headerMatch := matches[0]
				rest := line[strings.Index(line, headerMatch)+len(headerMatch):]
				activePart = stripComments(rest)
			}
		}

		if currentBlock != nil {
			// Check for heredoc start (supporting both <<- and <<) on activePart
			idx := strings.Index(activePart, "<<-")
			delimLen := 3
			if idx == -1 {
				idx = strings.Index(activePart, "<<")
				delimLen = 2
			}
			if idx != -1 {
				delimPart := strings.TrimSpace(activePart[idx+delimLen:])
				// Extract identifier, ignoring quotes if any
				delimPart = strings.Trim(delimPart, `"'`)
				if delimPart != "" {
					inHeredoc = true
					heredocDelim = delimPart

					// Count braces in the part of the line BEFORE the heredoc start
					beforeHeredoc := activePart[:idx]
					if strings.Contains(beforeHeredoc, "sensitive") {
						if sensitiveValRegex.MatchString(beforeHeredoc) {
							currentBlock.sensitive = true
						}
					}
					open, closedBraces := countBraces(beforeHeredoc)
					currentBlock.braces += open - closedBraces

					// If the block ended before the heredoc (unlikely), close it
					if currentBlock.braces <= 0 {
						if matchesKeywords(currentBlock.name, keywords) && !currentBlock.sensitive && !ignoreList.IsIgnored(filePath, currentBlock.name) {
							findings = append(findings, fmt.Sprintf("%s:%d: TF %s %q is not marked as sensitive", filePath, currentBlock.startLine, currentBlock.blockType, currentBlock.name))
						}
						currentBlock = nil
					}
					continue
				}
			}

			// Normal line processing (no heredoc start)
			if strings.Contains(activePart, "sensitive") {
				if sensitiveValRegex.MatchString(activePart) {
					currentBlock.sensitive = true
				}
			}

			open, closedBraces := countBraces(activePart)
			currentBlock.braces += open - closedBraces

			if currentBlock.braces <= 0 {
				// Block ended
				if matchesKeywords(currentBlock.name, keywords) && !currentBlock.sensitive && !ignoreList.IsIgnored(filePath, currentBlock.name) {
					findings = append(findings, fmt.Sprintf("%s:%d: TF %s %q is not marked as sensitive", filePath, currentBlock.startLine, currentBlock.blockType, currentBlock.name))
				}
				currentBlock = nil
			}
		}
	}

	return findings, nil
}

func main() {
	ignoreFilePath := flag.String("ignore", ".github/workflows/scripts/check_sensitive_ignore.cfg", "Path to sensitive ignore file")
	keywordsFilePath := flag.String("keywords", ".github/workflows/scripts/check_sensitive_keywords.cfg", "Path to sensitive keywords configuration file")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		args = []string{"rancher2", "examples"}
	}

	keywordsList, err := loadKeywords(*keywordsFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading keywords file: %v\n", err)
		os.Exit(1)
	}

	ignoreList, err := loadIgnoreList(*ignoreFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading ignore file: %v\n", err)
		os.Exit(1)
	}

	var allFindings []string
	fset := token.NewFileSet()

	for _, root := range args {
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error { // Use WalkDir for better performance
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}

			switch filepath.Ext(path) {
			case ".go":
				// Skip test files unless we want to analyze them too (usually not).
				if strings.HasSuffix(path, "_test.go") {
					return nil
				}

				node, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
				if err != nil {
					return fmt.Errorf("failed to parse Go file %s: %w", path, err)
				}

				findings := checkGoFile(fset, node, path, ignoreList, keywordsList)
				allFindings = append(allFindings, findings...)
			case ".tf":
				findings, err := checkTfFile(path, ignoreList, keywordsList)
				if err != nil {
					return fmt.Errorf("failed to analyze TF file %s: %w", path, err)
				}
				allFindings = append(allFindings, findings...)
			}

			return nil
		})

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error scanning %s: %v\n", root, err)
			os.Exit(1)
		}
	}

	if len(allFindings) > 0 {
		fmt.Printf("Found %d potentially sensitive fields not marked as sensitive:\n\n", len(allFindings))
		for _, finding := range allFindings {
			fmt.Println(finding)
		}
		fmt.Println("\n--------------------------------------------------------------------------------")
		fmt.Println("👉 ATTENTION DEVELOPER:")
		fmt.Println("If these fields are false positives and should NOT be marked as sensitive, please")
		fmt.Println("add them to the exemption configuration file:")
		fmt.Println("  .github/workflows/scripts/check_sensitive_ignore.cfg")
		fmt.Println("--------------------------------------------------------------------------------")
		os.Exit(1)
	}

	fmt.Println("No unmarked sensitive fields found. All looks good!")
}
