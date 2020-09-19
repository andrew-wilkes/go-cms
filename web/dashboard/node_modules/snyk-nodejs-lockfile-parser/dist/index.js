"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("source-map-support/register");
const fs = require("fs");
const path = require("path");
const parsers_1 = require("./parsers");
exports.Scope = parsers_1.Scope;
exports.LockfileType = parsers_1.LockfileType;
exports.getYarnWorkspaces = parsers_1.getYarnWorkspaces;
const package_lock_parser_1 = require("./parsers/package-lock-parser");
const yarn_lock_parse_1 = require("./parsers/yarn-lock-parse");
// import { Yarn2LockParser } from './parsers/yarn2-lock-parse';
// import getRuntimeVersion from './get-node-runtime-version';
const errors_1 = require("./errors");
exports.UnsupportedRuntimeError = errors_1.UnsupportedRuntimeError;
exports.InvalidUserInputError = errors_1.InvalidUserInputError;
exports.OutOfSyncError = errors_1.OutOfSyncError;
async function buildDepTree(manifestFileContents, lockFileContents, includeDev = false, lockfileType, strict = true, defaultManifestFileName = 'package.json') {
    if (!lockfileType) {
        lockfileType = parsers_1.LockfileType.npm;
    }
    else if (lockfileType === parsers_1.LockfileType.yarn) {
        lockfileType = getYarnLockfileType(lockFileContents);
    }
    let lockfileParser;
    switch (lockfileType) {
        case parsers_1.LockfileType.npm:
            lockfileParser = new package_lock_parser_1.PackageLockParser();
            break;
        case parsers_1.LockfileType.yarn:
            lockfileParser = new yarn_lock_parse_1.YarnLockParser();
            break;
        case parsers_1.LockfileType.yarn2:
            throw new errors_1.UnsupportedError('Yarn2 support has been temporarily removed to support Node.js versions 8.x.x');
        /**
         * Removing yarn 2 support as this breaks support for yarn with Node.js 8
         * See: https://github.com/snyk/snyk/issues/1270
         *
         * Uncomment following code once Snyk stops Node.js 8 support
         * // parsing yarn.lock is supported for Node.js v10 and higher
         * if (getRuntimeVersion() >= 10) {
         *  lockfileParser = new Yarn2LockParser();
         *  } else {
         *   throw new UnsupportedRuntimeError(
         *     'Parsing `yarn.lock` is not ' +
         *       'supported on Node.js version less than 10. Please upgrade your ' +
         *       'Node.js environment or use `package-lock.json`',
         *   );
         * }
         * break;
         */
        default:
            throw new errors_1.InvalidUserInputError('Unsupported lockfile type ' +
                `${lockfileType} provided. Only 'npm' or 'yarn' is currently ` +
                'supported.');
    }
    const manifestFile = parsers_1.parseManifestFile(manifestFileContents);
    if (!manifestFile.name) {
        manifestFile.name = path.isAbsolute(defaultManifestFileName)
            ? path.basename(defaultManifestFileName)
            : defaultManifestFileName;
    }
    const lockFile = lockfileParser.parseLockFile(lockFileContents);
    return lockfileParser.getDependencyTree(manifestFile, lockFile, includeDev, strict);
}
exports.buildDepTree = buildDepTree;
async function buildDepTreeFromFiles(root, manifestFilePath, lockFilePath, includeDev = false, strict = true) {
    if (!root || !manifestFilePath || !lockFilePath) {
        throw new Error('Missing required parameters for buildDepTreeFromFiles()');
    }
    const manifestFileFullPath = path.resolve(root, manifestFilePath);
    const lockFileFullPath = path.resolve(root, lockFilePath);
    if (!fs.existsSync(manifestFileFullPath)) {
        throw new errors_1.InvalidUserInputError('Target file package.json not found at ' +
            `location: ${manifestFileFullPath}`);
    }
    if (!fs.existsSync(lockFileFullPath)) {
        throw new errors_1.InvalidUserInputError('Lockfile not found at location: ' + lockFileFullPath);
    }
    const manifestFileContents = fs.readFileSync(manifestFileFullPath, 'utf-8');
    const lockFileContents = fs.readFileSync(lockFileFullPath, 'utf-8');
    let lockFileType;
    if (lockFilePath.endsWith('package-lock.json')) {
        lockFileType = parsers_1.LockfileType.npm;
    }
    else if (lockFilePath.endsWith('yarn.lock')) {
        lockFileType = getYarnLockfileType(lockFileContents, root, lockFilePath);
    }
    else {
        throw new errors_1.InvalidUserInputError(`Unknown lockfile ${lockFilePath}. ` +
            'Please provide either package-lock.json or yarn.lock.');
    }
    return await buildDepTree(manifestFileContents, lockFileContents, includeDev, lockFileType, strict, manifestFilePath);
}
exports.buildDepTreeFromFiles = buildDepTreeFromFiles;
function getYarnWorkspacesFromFiles(root, manifestFilePath) {
    if (!root || !manifestFilePath) {
        throw new Error('Missing required parameters for getYarnWorkspacesFromFiles()');
    }
    const manifestFileFullPath = path.resolve(root, manifestFilePath);
    if (!fs.existsSync(manifestFileFullPath)) {
        throw new errors_1.InvalidUserInputError('Target file package.json not found at ' +
            `location: ${manifestFileFullPath}`);
    }
    const manifestFileContents = fs.readFileSync(manifestFileFullPath, 'utf-8');
    return parsers_1.getYarnWorkspaces(manifestFileContents);
}
exports.getYarnWorkspacesFromFiles = getYarnWorkspacesFromFiles;
function getYarnLockfileType(lockFileContents, root, lockFilePath) {
    if (lockFileContents.includes('__metadata') ||
        (root &&
            lockFilePath &&
            fs.existsSync(path.resolve(root, lockFilePath.replace('yarn.lock', '.yarnrc.yml'))))) {
        return parsers_1.LockfileType.yarn2;
    }
    else {
        return parsers_1.LockfileType.yarn;
    }
}
//# sourceMappingURL=index.js.map