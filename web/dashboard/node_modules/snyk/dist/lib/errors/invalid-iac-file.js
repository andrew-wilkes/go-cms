"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.IllegalIacFileError = exports.NotSupportedIacFileError = void 0;
const chalk_1 = require("chalk");
const custom_error_1 = require("./custom-error");
function NotSupportedIacFileError(atLocations) {
    const locationsStr = atLocations.join(', ');
    const errorMsg = 'Not supported infrastructure as code target files in ' +
        locationsStr +
        '.\nPlease see our documentation for supported target files (currently we support Kubernetes files only): ' +
        chalk_1.default.underline('https://support.snyk.io/hc/en-us/articles/360006368877-Scan-and-fix-security-issues-in-your-Kubernetes-configuration-files') +
        ' and make sure you are in the right directory.';
    const error = new custom_error_1.CustomError(errorMsg);
    error.code = 422;
    error.userMessage = errorMsg;
    return error;
}
exports.NotSupportedIacFileError = NotSupportedIacFileError;
function IllegalIacFileError(atLocations) {
    const locationsStr = atLocations.join(', ');
    const errorMsg = 'Illegal infrastructure as code target file ' +
        locationsStr +
        '.\nPlease see our documentation for supported target files (currently we support Kubernetes files only): ' +
        chalk_1.default.underline('https://support.snyk.io/hc/en-us/articles/360006368877-Scan-and-fix-security-issues-in-your-Kubernetes-configuration-files') +
        ' and make sure you are in the right directory.';
    const error = new custom_error_1.CustomError(errorMsg);
    error.code = 422;
    error.userMessage = errorMsg;
    return error;
}
exports.IllegalIacFileError = IllegalIacFileError;
//# sourceMappingURL=invalid-iac-file.js.map