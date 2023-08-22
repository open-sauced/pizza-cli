#!/usr/bin/env node

const { spawnSync } = require("child_process");
const path = require("path");
const { name } = require("../package.json");

const command_args = process.argv.slice(2);
const exeName = ["win32", "cygwin"].includes(process.platform)
  ? `${name}.exe`
  : name;
const binPath = path.join(__dirname, exeName);
const child = spawnSync(binPath, command_args, { stdio: "inherit" });
process.exit(child.status ?? 1);
