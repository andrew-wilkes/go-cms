"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.testDependencies = exports.testEcosystem = exports.getEcosystem = exports.getPlugin = void 0;
const cppPlugin = require("snyk-cpp-plugin");
const snyk = require("./index");
const config = require("./config");
const is_ci_1 = require("./is-ci");
const promise_1 = require("./request/promise");
const types_1 = require("../cli/commands/types");
const spinner = require("../lib/spinner");
const EcosystemPlugins = {
    cpp: cppPlugin,
};
function getPlugin(ecosystem) {
    return EcosystemPlugins[ecosystem];
}
exports.getPlugin = getPlugin;
function getEcosystem(options) {
    if (options.source) {
        return 'cpp';
    }
    return null;
}
exports.getEcosystem = getEcosystem;
async function testEcosystem(ecosystem, paths, options) {
    const plugin = getPlugin(ecosystem);
    const scanResultsByPath = {};
    for (const path of paths) {
        options.path = path;
        const results = await plugin.scan(options);
        scanResultsByPath[path] = results;
    }
    const [testResults, errors] = await testDependencies(scanResultsByPath);
    const stringifiedData = JSON.stringify(testResults, null, 2);
    if (options.json) {
        return types_1.TestCommandResult.createJsonTestCommandResult(stringifiedData);
    }
    const emptyResults = [];
    const scanResults = emptyResults.concat(...Object.values(scanResultsByPath));
    const readableResult = await plugin.display(scanResults, testResults, errors, options);
    return types_1.TestCommandResult.createHumanReadableTestCommandResult(readableResult, stringifiedData);
}
exports.testEcosystem = testEcosystem;
async function testDependencies(scans) {
    const results = [];
    const errors = [];
    for (const [path, scanResults] of Object.entries(scans)) {
        await spinner(`Testing dependencies in ${path}`);
        for (const scanResult of scanResults) {
            const payload = {
                method: 'POST',
                url: `${config.API}/test-dependencies`,
                json: true,
                headers: {
                    'x-is-ci': is_ci_1.isCI(),
                    authorization: 'token ' + snyk.api,
                },
                body: {
                    artifacts: scanResult.artifacts,
                    meta: {},
                },
            };
            try {
                const response = await promise_1.makeRequest(payload);
                results.push(response);
            }
            catch (error) {
                if (error.code >= 400 && error.code < 500) {
                    throw new Error(error.message);
                }
                errors.push('Could not test dependencies in ' + path);
            }
        }
    }
    spinner.clearAll();
    return [results, errors];
}
exports.testDependencies = testDependencies;
//# sourceMappingURL=ecosystems.js.map