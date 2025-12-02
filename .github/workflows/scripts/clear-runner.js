import { execSync } from 'child_process';
export default async ({ core }) => {
  const pathsToRemove = [
    '/usr/lib/jvm',
    '/usr/share/dotnet',
    '/usr/share/swift',
    '/usr/local/.ghcup',
    '/usr/local/julia*',
    '/usr/local/lib/android',
    '/usr/local/share/chromium',
    '/opt/microsoft',
    '/opt/google',
    '/opt/az',
    '/usr/local/share/powershell',
    '/opt/hostedtoolcache'
  ];

  core.info('Disk space before cleanup:');
  try {
    execSync(`df -h`);
  } catch (error) {
    core.setFailed(`Failed running df to see disk space: ${error}`);
  }

  // Iterate over paths and remove them
  for (const path of pathsToRemove) {
    core.info(`Removing ${path}...`);
    try {
      // We use 'bash -c' to ensure wildcards (like julia*) are expanded correctly
      execSync(`sudo bash -c "rm -rf ${path}"`);
    } catch (error) {
      core.setFailed(`Failed to remove ${path}: ${error}`);
    }
  }

  core.info('Pruning Docker...');
  try {
    execSync(`docker system prune -af`);
    execSync(`docker builder prune -af`);
    execSync(`docker image prune -af`);
    execSync(`docker volume prune -af`);
  } catch(error) {
    core.setFailed(`Failed pruning Docker: ${error}`);
  }

  core.info('Disk space after cleanup:');
  try {
    execSync(`df -h`);
  } catch (error) {
    core.setFailed(`Failed running df to see disk space: ${error}`);
  }
};
