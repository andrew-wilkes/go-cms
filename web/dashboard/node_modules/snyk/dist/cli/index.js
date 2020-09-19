#!/usr/bin/env node
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
require("source-map-support/register");
const Debug = require("debug");
const pathLib = require("path");
// assert supported node runtime version
const runtime = require("./runtime");
// require analytics as soon as possible to start measuring execution time
const analytics = require("../lib/analytics");
const alerts = require("../lib/alerts");
const sln = require("../lib/sln");
const args_1 = require("./args");
const copy_1 = require("./copy");
const spinner = require("../lib/spinner");
const errors = require("../lib/errors/legacy-errors");
const ansiEscapes = require("ansi-escapes");
const detect_1 = require("../lib/detect");
const updater_1 = require("../lib/updater");
const errors_1 = require("../lib/errors");
const strip_ansi_1 = require("strip-ansi");
const exclude_flag_invalid_input_1 = require("../lib/errors/exclude-flag-invalid-input");
const modes_1 = require("./modes");
const json_file_output_bad_input_error_1 = require("../lib/errors/json-file-output-bad-input-error");
const json_file_output_1 = require("../lib/json-file-output");
const empty_sarif_output_error_1 = require("../lib/errors/empty-sarif-output-error");
const debug = Debug('snyk');
const EXIT_CODES = {
    VULNS_FOUND: 1,
    ERROR: 2,
    NO_SUPPORTED_MANIFESTS_FOUND: 3,
};
async function runCommand(args) {
    const commandResult = await args.method(...args.options._);
    const res = analytics({
        args: args.options._,
        command: args.command,
        org: args.options.org,
    });
    if (!commandResult) {
        return;
    }
    const result = commandResult.toString();
    if (result && !args.options.quiet) {
        if (args.options.copy) {
            copy_1.copy(result);
            console.log('Result copied to clipboard');
        }
        else {
            console.log(result);
        }
    }
    // also save the json (in error.json) to file if option is set
    if (args.command === 'test') {
        const jsonResults = commandResult.getJsonResult();
        saveResultsToFile(args.options, 'json', jsonResults);
        const sarifResults = commandResult.getSarifResult();
        saveResultsToFile(args.options, 'sarif', sarifResults);
    }
    return res;
}
async function handleError(args, error) {
    var _a;
    spinner.clearAll();
    let command = 'bad-command';
    let exitCode = EXIT_CODES.ERROR;
    const noSupportedManifestsFound = (_a = error.message) === null || _a === void 0 ? void 0 : _a.includes('Could not detect supported target files in');
    if (noSupportedManifestsFound) {
        exitCode = EXIT_CODES.NO_SUPPORTED_MANIFESTS_FOUND;
    }
    const vulnsFound = error.code === 'VULNS';
    if (vulnsFound) {
        // this isn't a bad command, so we won't record it as such
        command = args.command;
        exitCode = EXIT_CODES.VULNS_FOUND;
    }
    if (args.options.debug && !args.options.json) {
        const output = vulnsFound ? error.message : error.stack;
        console.log(output);
    }
    else if (args.options.json) {
        console.log(strip_ansi_1.default(error.json || error.stack));
    }
    else {
        if (!args.options.quiet) {
            const result = errors.message(error);
            if (args.options.copy) {
                copy_1.copy(result);
                console.log('Result copied to clipboard');
            }
            else {
                if (`${error.code}`.indexOf('AUTH_') === 0) {
                    // remove the last few lines
                    const erase = ansiEscapes.eraseLines(4);
                    process.stdout.write(erase);
                }
                console.log(result);
            }
        }
    }
    saveResultsToFile(args.options, 'json', error.jsonStringifiedResults);
    saveResultsToFile(args.options, 'sarif', error.sarifStringifiedResults);
    const analyticsError = vulnsFound
        ? {
            stack: error.jsonNoVulns,
            code: error.code,
            message: 'Vulnerabilities found',
        }
        : {
            stack: error.stack,
            code: error.code,
            message: error.message,
        };
    if (!vulnsFound && !error.stack) {
        // log errors that are not error objects
        analytics.add('error', true);
        analytics.add('command', args.command);
    }
    else {
        analytics.add('error-message', analyticsError.message);
        // Note that error.stack would also contain the error message
        // (see https://nodejs.org/api/errors.html#errors_error_stack)
        analytics.add('error', analyticsError.stack);
        analytics.add('error-code', error.code);
        analytics.add('command', args.command);
    }
    const res = analytics({
        args: args.options._,
        command,
        org: args.options.org,
    });
    return { res, exitCode };
}
function getFullPath(filepathFragment) {
    if (pathLib.isAbsolute(filepathFragment)) {
        return filepathFragment;
    }
    else {
        const fullPath = pathLib.join(process.cwd(), filepathFragment);
        return fullPath;
    }
}
function saveJsonResultsToFile(stringifiedJson, jsonOutputFile) {
    if (!jsonOutputFile) {
        console.error('empty jsonOutputFile');
        return;
    }
    if (jsonOutputFile.constructor.name !== String.name) {
        console.error('--json-output-file should be a filename path');
        return;
    }
    // create the directory if it doesn't exist
    const dirPath = pathLib.dirname(jsonOutputFile);
    const createDirSuccess = json_file_output_1.createDirectory(dirPath);
    if (createDirSuccess) {
        json_file_output_1.writeContentsToFileSwallowingErrors(jsonOutputFile, stringifiedJson);
    }
}
function checkRuntime() {
    if (!runtime.isSupported(process.versions.node)) {
        console.error(`${process.versions.node} is an unsupported nodejs ` +
            `runtime! Supported runtime range is '${runtime.supportedRange}'`);
        console.error('Please upgrade your nodejs runtime version and try again.');
        process.exit(EXIT_CODES.ERROR);
    }
}
// Throw error if user specifies package file name as part of path,
// and if user specifies multiple paths and used project-name option.
function checkPaths(args) {
    let count = 0;
    for (const path of args.options._) {
        if (typeof path === 'string' && detect_1.isPathToPackageFile(path)) {
            throw errors_1.MissingTargetFileError(path);
        }
        else if (typeof path === 'string') {
            if (++count > 1 && args.options['project-name']) {
                throw new errors_1.UnsupportedOptionCombinationError([
                    'multiple paths',
                    'project-name',
                ]);
            }
        }
    }
}
async function main() {
    updater_1.updateCheck();
    checkRuntime();
    const args = args_1.args(process.argv);
    let res;
    let failed = false;
    let exitCode = EXIT_CODES.ERROR;
    try {
        modes_1.modeValidation(args);
        // TODO: fix this, we do transformation to options and teh type doesn't reflect it
        validateUnsupportedOptionCombinations(args.options);
        if (args.options['app-vulns'] && args.options['json']) {
            throw new errors_1.UnsupportedOptionCombinationError([
                'Application vulnerabilities is currently not supported with JSON output. ' +
                    'Please try using —app-vulns only to get application vulnerabilities, or ' +
                    '—json only to get your image vulnerabilties, excluding the application ones.',
            ]);
        }
        if (args.options.file &&
            typeof args.options.file === 'string' &&
            args.options.file.match(/\.sln$/)) {
            if (args.options['project-name']) {
                throw new errors_1.UnsupportedOptionCombinationError([
                    'file=*.sln',
                    'project-name',
                ]);
            }
            sln.updateArgs(args);
        }
        else if (typeof args.options.file === 'boolean') {
            throw new errors_1.FileFlagBadInputError();
        }
        validateUnsupportedSarifCombinations(args);
        validateOutputFile(args.options, 'json', new json_file_output_bad_input_error_1.JsonFileOutputBadInputError());
        validateOutputFile(args.options, 'sarif', new empty_sarif_output_error_1.SarifFileOutputEmptyError());
        checkPaths(args);
        res = await runCommand(args);
    }
    catch (error) {
        failed = true;
        const response = await handleError(args, error);
        res = response.res;
        exitCode = response.exitCode;
    }
    if (!args.options.json) {
        console.log(alerts.displayAlerts());
    }
    if (!process.env.TAP && failed) {
        debug('Exit code: ' + exitCode);
        process.exitCode = exitCode;
    }
    return res;
}
const cli = main().catch((e) => {
    console.error('Something unexpected went wrong: ', e.stack);
    console.error('Exit code: ' + EXIT_CODES.ERROR);
    process.exit(EXIT_CODES.ERROR);
});
if (module.parent) {
    // eslint-disable-next-line id-blacklist
    module.exports = cli;
}
function validateUnsupportedOptionCombinations(options) {
    const unsupportedAllProjectsCombinations = {
        'project-name': 'project-name',
        file: 'file',
        yarnWorkspaces: 'yarn-workspaces',
        packageManager: 'package-manager',
        docker: 'docker',
        allSubProjects: 'all-sub-projects',
    };
    const unsupportedYarnWorkspacesCombinations = {
        'project-name': 'project-name',
        file: 'file',
        packageManager: 'package-manager',
        docker: 'docker',
        allSubProjects: 'all-sub-projects',
    };
    if (options.scanAllUnmanaged && options.file) {
        throw new errors_1.UnsupportedOptionCombinationError(['file', 'scan-all-unmanaged']);
    }
    if (options.allProjects) {
        for (const option in unsupportedAllProjectsCombinations) {
            if (options[option]) {
                throw new errors_1.UnsupportedOptionCombinationError([
                    unsupportedAllProjectsCombinations[option],
                    'all-projects',
                ]);
            }
        }
    }
    if (options.yarnWorkspaces) {
        for (const option in unsupportedYarnWorkspacesCombinations) {
            if (options[option]) {
                throw new errors_1.UnsupportedOptionCombinationError([
                    unsupportedAllProjectsCombinations[option],
                    'yarn-workspaces',
                ]);
            }
        }
    }
    if (options.exclude) {
        if (!(options.allProjects || options.yarnWorkspaces)) {
            throw new errors_1.OptionMissingErrorError('--exclude', [
                '--yarn-workspaces',
                '--all-projects',
            ]);
        }
        if (typeof options.exclude !== 'string') {
            throw new errors_1.ExcludeFlagBadInputError();
        }
        if (options.exclude.indexOf(pathLib.sep) > -1) {
            throw new exclude_flag_invalid_input_1.ExcludeFlagInvalidInputError();
        }
    }
}
function validateUnsupportedSarifCombinations(args) {
    if (args.options['json-file-output'] && args.command !== 'test') {
        throw new errors_1.UnsupportedOptionCombinationError([
            args.command,
            'json-file-output',
        ]);
    }
    if (args.options['sarif'] && args.command !== 'test') {
        throw new errors_1.UnsupportedOptionCombinationError([args.command, 'sarif']);
    }
    if (args.options['sarif'] && args.options['json']) {
        throw new errors_1.UnsupportedOptionCombinationError([
            args.command,
            'sarif',
            'json',
        ]);
    }
    if (args.options['sarif-file-output'] && args.command !== 'test') {
        throw new errors_1.UnsupportedOptionCombinationError([
            args.command,
            'sarif-file-output',
        ]);
    }
    if (args.options['sarif'] &&
        args.options['docker'] &&
        !args.options['file']) {
        throw new errors_1.OptionMissingErrorError('sarif', ['--file']);
    }
    if (args.options['sarif-file-output'] &&
        args.options['docker'] &&
        !args.options['file']) {
        throw new errors_1.OptionMissingErrorError('sarif-file-output', ['--file']);
    }
}
function saveResultsToFile(options, outputType, jsonResults) {
    const flag = `${outputType}-file-output`;
    const outputFile = options[flag];
    if (outputFile && jsonResults) {
        const outputFileStr = outputFile;
        const fullOutputFilePath = getFullPath(outputFileStr);
        saveJsonResultsToFile(strip_ansi_1.default(jsonResults), fullOutputFilePath);
    }
}
function validateOutputFile(options, outputType, error) {
    const fileOutputValue = options[`${outputType}-file-output`];
    if (fileOutputValue === undefined) {
        return;
    }
    if (!fileOutputValue || typeof fileOutputValue !== 'string') {
        throw error;
    }
    // On Windows, seems like quotes get passed in
    if (fileOutputValue === "''" || fileOutputValue === '""') {
        throw error;
    }
}
//# sourceMappingURL=index.js.map