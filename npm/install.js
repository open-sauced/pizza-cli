const fs = require("fs");
const https = require("https");
const path = require("path");

const packageJson = require("./package.json");
const { version, name, repository } = packageJson;

async function install() {
  const downloadURL = getDownloadURL();

  try {
    const binDir = path.join(__dirname, "bin");
    const exeName = ["win32", "cygwin"].includes(process.platform)
      ? `${name}.exe`
      : name;
    const outputPath = path.join(binDir, exeName);

    await downloadBinary(downloadURL, outputPath);

    fs.chmodSync(outputPath, 0o755);
  } catch (error) {
    console.error("Installation failed:", error.message);
  }
}

function getDownloadURL() {
  let goOS, arch;

  switch (process.platform) {
    case "win32":
    case "cygwin":
      goOS = "windows";
      break;
    case "darwin":
      goOS = "darwin";
      break;
    case "linux":
      goOS = "linux";
      break;
    default:
      throw new Error(`Unsupported OS: ${process.platform}`);
  }

  switch (process.arch) {
    case "x64":
      arch = "amd64";
      break;
    case "arm64":
      arch = "arm64";
      break;
    default:
      throw new Error(`Unsupported architecture: ${process.arch}`);
  }

  return `${repository}/releases/download/v${version}/${name}-${goOS}-${arch}`;
}

const downloadBinary = (url, outputPath) => {
  return new Promise((resolve, reject) => {
    https
      .get(url, (response) => {
        if (response.statusCode === 302) {
          resolve(downloadBinary(response.headers.location, outputPath));
        } else if (response.statusCode === 200) {
          const file = fs.createWriteStream(outputPath);
          response.pipe(file);
          file.on("finish", () => {
            file.close(resolve);
          });
        } else {
          reject(
            new Error(
              `Failed to download ${name}. Status code: ${response.statusCode}`
            )
          );
        }
      })
      .on("error", reject);
  });
};

void install();
