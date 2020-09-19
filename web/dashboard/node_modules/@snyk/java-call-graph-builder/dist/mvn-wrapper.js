"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getClassPathFromMvn = exports.mergeMvnClassPaths = exports.parseMvnExecCommandOutput = exports.parseMvnDependencyPluginCommandOutput = void 0;
const tslib_1 = require("tslib");
require("source-map-support/register");
const sub_process_1 = require("./sub-process");
const errors_1 = require("./errors");
function getMvnCommandArgsForMvnExec(targetPath) {
    return [
        '-q',
        'exec:exec',
        '-Dexec.classpathScope="compile"',
        '-Dexec.executable="echo"',
        '-Dexec.args="%classpath"',
        '-f',
        targetPath,
    ];
}
function getMvnCommandArgsForDependencyPlugin(targetPath) {
    return ['dependency:build-classpath', '-f', targetPath];
}
function parseMvnDependencyPluginCommandOutput(mvnCommandOutput) {
    const outputLines = mvnCommandOutput.split('\n');
    const mvnClassPaths = [];
    let startIndex = 0;
    let i = outputLines.indexOf('[INFO] Dependencies classpath:', startIndex);
    while (i > -1) {
        if (outputLines[i + 1] !== '') {
            mvnClassPaths.push(outputLines[i + 1]);
        }
        startIndex = i + 2;
        i = outputLines.indexOf('[INFO] Dependencies classpath:', startIndex);
    }
    return mvnClassPaths;
}
exports.parseMvnDependencyPluginCommandOutput = parseMvnDependencyPluginCommandOutput;
function parseMvnExecCommandOutput(mvnCommandOutput) {
    return mvnCommandOutput.trim().split('\n');
}
exports.parseMvnExecCommandOutput = parseMvnExecCommandOutput;
function mergeMvnClassPaths(classPaths) {
    // this magic joins all items in array with :, splits result by : again
    // makes Set (to uniq items), create Array from it and join it by : to have
    // proper path like format
    return Array.from(new Set(classPaths.join(':').split(':'))).join(':');
}
exports.mergeMvnClassPaths = mergeMvnClassPaths;
function getClassPathFromMvn(targetPath) {
    return tslib_1.__awaiter(this, void 0, void 0, function* () {
        let classPaths = [];
        let args = [];
        try {
            try {
                // there are two ways of getting classpath - either from maven plugin or by exec command
                // try `mvn exec` for classpath
                args = getMvnCommandArgsForMvnExec(targetPath);
                const output = yield sub_process_1.execute('mvn', args, { cwd: targetPath });
                classPaths = parseMvnExecCommandOutput(output);
            }
            catch (e) {
                // if it fails, try mvn dependency:build-classpath
                // TODO send error message for further analysis
                args = getMvnCommandArgsForDependencyPlugin(targetPath);
                const output = yield sub_process_1.execute('mvn', args, { cwd: targetPath });
                classPaths = parseMvnDependencyPluginCommandOutput(output);
            }
        }
        catch (e) {
            throw new errors_1.ClassPathGenerationError(e);
        }
        if (classPaths.length === 0) {
            throw new errors_1.EmptyClassPathError(`mvn ${args.join(' ')}`);
        }
        return mergeMvnClassPaths(classPaths);
    });
}
exports.getClassPathFromMvn = getClassPathFromMvn;
//# sourceMappingURL=mvn-wrapper.js.map