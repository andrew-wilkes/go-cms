"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubprocessError = exports.SubprocessTimeoutError = exports.MissingTargetFolderError = exports.EmptyClassPathError = exports.ClassPathGenerationError = exports.CallGraphGenerationError = void 0;
class CallGraphGenerationError extends Error {
    constructor(msg, innerError) {
        super(msg);
        Object.setPrototypeOf(this, CallGraphGenerationError.prototype);
        this.innerError = innerError;
    }
}
exports.CallGraphGenerationError = CallGraphGenerationError;
class ClassPathGenerationError extends Error {
    constructor(innerError) {
        super('Class path generation error');
        this.userMessage = "Could not determine the project's class path. Please contact our support or submit an issue at https://github.com/snyk/java-call-graph-builder/issues. Re-running the command with the `-d` flag will provide useful information for the support engineers.";
        Object.setPrototypeOf(this, ClassPathGenerationError.prototype);
        this.innerError = innerError;
    }
}
exports.ClassPathGenerationError = ClassPathGenerationError;
class EmptyClassPathError extends Error {
    constructor(command) {
        super(`The command "${command}" returned an empty class path`);
        this.userMessage = 'The class path for the project is empty. Please contact our support or submit an issue at https://github.com/snyk/java-call-graph-builder/issues. Re-running the command with the `-d` flag will provide useful information for the support engineers.';
        Object.setPrototypeOf(this, EmptyClassPathError.prototype);
    }
}
exports.EmptyClassPathError = EmptyClassPathError;
class MissingTargetFolderError extends Error {
    constructor(targetPath) {
        super(`Could not find the target folder starting in "${targetPath}"`);
        this.userMessage = "Could not find the project's target folder. Please compile your code by running `mvn compile` and try again.";
        Object.setPrototypeOf(this, MissingTargetFolderError.prototype);
    }
}
exports.MissingTargetFolderError = MissingTargetFolderError;
class SubprocessTimeoutError extends Error {
    constructor(command, args, timeout) {
        super(`The command "${command} ${args}" timed out after ${timeout / 1000}s`);
        this.userMessage = 'Scanning for reachable vulnerabilities took too long. Please use the --reachable-timeout flag to increase the timeout for finding reachable vulnerabilities.';
        Object.setPrototypeOf(this, SubprocessTimeoutError.prototype);
    }
}
exports.SubprocessTimeoutError = SubprocessTimeoutError;
class SubprocessError extends Error {
    constructor(command, args, exitCode) {
        super(`The command "${command} ${args}" exited with code ${exitCode}`);
        Object.setPrototypeOf(this, SubprocessError.prototype);
    }
}
exports.SubprocessError = SubprocessError;
//# sourceMappingURL=errors.js.map