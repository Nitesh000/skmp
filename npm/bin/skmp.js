#!/usr/bin/env node
const { execFileSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const bin = path.join(__dirname, process.platform === 'win32' ? 'skmp-native.exe' : 'skmp-native');

if (!fs.existsSync(bin)) {
  console.error(`skmp binary not found at ${bin}`);
  console.error('Please run "npm run postinstall" to download it.');
  process.exit(1);
}

try {
  execFileSync(bin, process.argv.slice(2), { stdio: 'inherit' });
} catch (err) {
  process.exit(err.status || 1);
}
