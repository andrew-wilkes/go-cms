"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.runtimeMetrics = exports.getClassGraphGradle = exports.getCallGraphMvn = void 0;
const tslib_1 = require("tslib");
require("source-map-support/register");
const mvn_wrapper_1 = require("./mvn-wrapper");
const gradle_wrapper_1 = require("./gradle-wrapper");
const java_wrapper_1 = require("./java-wrapper");
const metrics_1 = require("./metrics");
const errors_1 = require("./errors");
function getCallGraphMvn(targetPath, timeout) {
    return tslib_1.__awaiter(this, void 0, void 0, function* () {
        try {
            const classPath = yield metrics_1.timeIt('getMvnClassPath', () => mvn_wrapper_1.getClassPathFromMvn(targetPath));
            return yield metrics_1.timeIt('getCallGraph', () => java_wrapper_1.getCallGraph(classPath, targetPath, timeout));
        }
        catch (e) {
            throw new errors_1.CallGraphGenerationError(e.userMessage ||
                'Failed to scan for reachable vulnerabilities. Please contact our support or submit an issue at https://github.com/snyk/java-call-graph-builder/issues. Re-running the command with the `-d` flag will provide useful information for the support engineers.', e);
        }
    });
}
exports.getCallGraphMvn = getCallGraphMvn;
function getClassGraphGradle(targetPath, timeout) {
    return tslib_1.__awaiter(this, void 0, void 0, function* () {
        const classPath = yield metrics_1.timeIt('getGradleClassPath', () => gradle_wrapper_1.getClassPathFromGradle(targetPath));
        return yield metrics_1.timeIt('getCallGraph', () => java_wrapper_1.getCallGraph(classPath, targetPath, timeout));
    });
}
exports.getClassGraphGradle = getClassGraphGradle;
function runtimeMetrics() {
    return metrics_1.getMetrics();
}
exports.runtimeMetrics = runtimeMetrics;
//# sourceMappingURL=index.js.map