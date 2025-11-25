async () => {
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
  await exec.exec('df', ['-h']);

  // Iterate over paths and remove them
  for (const path of pathsToRemove) {
    core.info(`Removing ${path}...`);
    // We use 'bash -c' to ensure wildcards (like julia*) are expanded correctly
    // We use ignoreReturnCode: true so the build doesn't fail if a folder is already missing
    await exec.exec('sudo', ['bash', '-c', `rm -rf ${path}`], { ignoreReturnCode: true });
  }

  core.info('Pruning Docker...');
  await exec.exec('docker', ['system', 'prune', '-af'], { ignoreReturnCode: true });
  await exec.exec('docker', ['builder', 'prune', '-af'], { ignoreReturnCode: true });

  core.info('Disk space after cleanup:');
  await exec.exec('df', ['-h']);
};
