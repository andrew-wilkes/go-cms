"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getWorkspacesMap = exports.processYarnWorkspaces = void 0;
const baseDebug = require("debug");
const pathUtil = require("path");
const _ = require("lodash");
const debug = baseDebug('snyk:yarn-workspaces');
const lockFileParser = require("snyk-nodejs-lockfile-parser");
const path = require("path");
const errors_1 = require("../../errors");
const get_file_contents_1 = require("../../get-file-contents");
async function processYarnWorkspaces(root, settings, targetFiles) {
    // the order of folders is important
    // must have the root level most folders at the top
    const yarnTargetFiles = _(targetFiles)
        .map((p) => (Object.assign({ path: p }, pathUtil.parse(p))))
        .filter((res) => ['package.json'].includes(res.base))
        .sortBy('dir')
        .groupBy('dir')
        .value();
    debug(`processing Yarn workspaces (${targetFiles.length})`);
    if (Object.keys(yarnTargetFiles).length === 0) {
        throw errors_1.NoSupportedManifestsFoundError([root]);
    }
    let yarnWorkspacesMap = {};
    const yarnWorkspacesFilesMap = {};
    let isYarnWorkspacePackage = false;
    const result = {
        plugin: {
            name: 'snyk-nodejs-yarn-workspaces',
            runtime: process.version,
        },
        scannedProjects: [],
    };
    // the folders must be ordered highest first
    for (const directory of Object.keys(yarnTargetFiles)) {
        const packageJsonFileName = pathUtil.join(directory, 'package.json');
        const packageJson = get_file_contents_1.getFileContents(root, packageJsonFileName);
        yarnWorkspacesMap = Object.assign(Object.assign({}, yarnWorkspacesMap), getWorkspacesMap(packageJson));
        for (const workspaceRoot of Object.keys(yarnWorkspacesMap)) {
            const workspaces = yarnWorkspacesMap[workspaceRoot].workspaces || [];
            const match = workspaces
                .map((pattern) => {
                return packageJsonFileName.includes(pattern.replace(/\*/, ''));
            })
                .filter(Boolean);
            if (match) {
                yarnWorkspacesFilesMap[packageJsonFileName] = {
                    root: workspaceRoot,
                };
                isYarnWorkspacePackage = true;
            }
        }
        if (isYarnWorkspacePackage) {
            const rootDir = path.dirname(yarnWorkspacesFilesMap[packageJsonFileName].root);
            const rootYarnLockfileName = path.join(rootDir, 'yarn.lock');
            const yarnLock = await get_file_contents_1.getFileContents(root, rootYarnLockfileName);
            const res = await lockFileParser.buildDepTree(packageJson.content, yarnLock.content, settings.dev, lockFileParser.LockfileType.yarn, settings.strictOutOfSync !== false);
            const project = {
                packageManager: 'yarn',
                targetFile: path.relative(root, packageJson.fileName),
                depTree: res,
                plugin: {
                    name: 'snyk-nodejs-lockfile-parser',
                    runtime: process.version,
                },
            };
            result.scannedProjects.push(project);
        }
    }
    return result;
}
exports.processYarnWorkspaces = processYarnWorkspaces;
function getWorkspacesMap(file) {
    const yarnWorkspacesMap = {};
    if (!file) {
        return yarnWorkspacesMap;
    }
    try {
        const rootFileWorkspacesDefinitions = lockFileParser.getYarnWorkspaces(file.content);
        if (rootFileWorkspacesDefinitions && rootFileWorkspacesDefinitions.length) {
            yarnWorkspacesMap[file.fileName] = {
                workspaces: rootFileWorkspacesDefinitions,
            };
        }
    }
    catch (e) {
        debug('Failed to process a workspace', e.message);
    }
    return yarnWorkspacesMap;
}
exports.getWorkspacesMap = getWorkspacesMap;
//# sourceMappingURL=yarn-workspaces-parser.js.map