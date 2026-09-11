const fs = require('fs');
const path = require('path');
const https = require('https');

const VERSION = require('../package.json').version;
const REPO = 'Nitesh000/skmp';

const osMap = {
  darwin: 'darwin',
  linux: 'linux',
  win32: 'windows'
};

const archMap = {
  x64: 'amd64',
  arm64: 'arm64'
};

const os = osMap[process.platform];
const arch = archMap[process.arch];

if (!os || !arch) {
  console.error(`Unsupported platform: ${process.platform} ${process.arch}`);
  process.exit(1);
}

const ext = os === 'windows' ? '.exe' : '';
const binName = `skmp-${os}-${arch}${ext}`;
const url = `https://github.com/${REPO}/releases/download/v${VERSION}/${binName}`;
const dest = path.join(__dirname, '..', 'bin', os === 'windows' ? 'skmp-native.exe' : 'skmp-native');

console.log(`Downloading skmp binary for ${os}-${arch} from ${url}...`);

function download(url, dest) {
  return new Promise((resolve, reject) => {
    https.get(url, (res) => {
      if (res.statusCode === 301 || res.statusCode === 302) {
        return download(res.headers.location, dest).then(resolve).catch(reject);
      }
      
      if (res.statusCode !== 200) {
        return reject(new Error(`Failed to download binary: HTTP ${res.statusCode}`));
      }
      
      const file = fs.createWriteStream(dest);
      res.pipe(file);
      
      file.on('finish', () => {
        file.close();
        if (os !== 'windows') {
          fs.chmodSync(dest, 0o755);
        }
        resolve();
      });
    }).on('error', (err) => {
      fs.unlink(dest, () => reject(err));
    });
  });
}

download(url, dest)
  .then(() => {
    console.log('skmp binary successfully installed.');
  })
  .catch((err) => {
    console.error(err.message);
    process.exit(1);
  });
