"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getClassPathFromGradle = exports.getGradleCommand = exports.getGradleCommandArgs = void 0;
const tslib_1 = require("tslib");
require("source-map-support/register");
const sub_process_1 = require("./sub-process");
const path = require("path");
const fs = require("fs");
const errors_1 = require("./errors");
function getGradleCommandArgs(targetPath) {
    const gradleArgs = [
        'printClasspath',
        '-I',
        path.join(__dirname, ...'../bin/init.gradle'.split('/')),
        '-q',
    ];
    if (targetPath) {
        gradleArgs.push('-p', targetPath);
    }
    return gradleArgs;
}
exports.getGradleCommandArgs = getGradleCommandArgs;
function getGradleCommand(targetPath) {
    const pathToWrapper = path.resolve(targetPath || '', '.', 'gradlew');
    if (fs.existsSync(pathToWrapper)) {
        return pathToWrapper;
    }
    return 'gradle';
}
exports.getGradleCommand = getGradleCommand;
function getClassPathFromGradle(targetPath) {
    return tslib_1.__awaiter(this, void 0, void 0, function* () {
        const cmd = getGradleCommand(targetPath);
        const args = getGradleCommandArgs(targetPath);
        try {
            const output = yield sub_process_1.execute(cmd, args, { cwd: targetPath });
            return output.trim();
        }
        catch (e) {
            throw new errors_1.ClassPathGenerationError(e);
        }
    });
}
exports.getClassPathFromGradle = getClassPathFromGradle;
//# sourceMappingURL=gradle-wrapper.js.map