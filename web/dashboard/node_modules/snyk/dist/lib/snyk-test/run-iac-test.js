"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.assembleIacLocalPayloads = exports.parseIacTestResult = void 0;
const fs = require("fs");
const path = require("path");
const pathUtil = require("path");
const snyk = require("..");
const is_ci_1 = require("../is-ci");
const common = require("./common");
const config = require("../config");
const legacy_1 = require("./legacy");
const pathLib = require("path");
async function parseIacTestResult(res, targetFile, projectName, severityThreshold) {
    const meta = res.meta || {};
    severityThreshold =
        severityThreshold === legacy_1.SEVERITY.LOW ? undefined : severityThreshold;
    return Object.assign(Object.assign({}, res), { vulnerabilities: [], dependencyCount: 0, licensesPolicy: null, ignoreSettings: null, targetFile,
        projectName, org: meta.org, policy: meta.policy, isPrivate: !meta.isPublic, severityThreshold });
}
exports.parseIacTestResult = parseIacTestResult;
async function assembleIacLocalPayloads(root, options) {
    const payloads = [];
    // Forcing options.path to be a string as pathUtil requires is to be stringified
    const targetFile = pathLib.resolve(root, '.');
    const targetFileRelativePath = targetFile
        ? pathUtil.join(pathUtil.resolve(`${options.path}`), targetFile)
        : '';
    const fileContent = fs.readFileSync(targetFile, 'utf8');
    const body = {
        data: {
            fileContent,
            fileType: 'yaml',
        },
        targetFile: root,
        type: 'k8sconfig',
        //TODO(orka): future - support policy
        policy: '',
        targetFileRelativePath: `${targetFileRelativePath}`,
        originalProjectName: path.basename(path.dirname(targetFile)),
        projectNameOverride: options.projectName,
    };
    const payload = {
        method: 'POST',
        url: config.API + (options.vulnEndpoint || '/test-iac'),
        json: true,
        headers: {
            'x-is-ci': is_ci_1.isCI(),
            authorization: 'token ' + snyk.api,
        },
        qs: common.assembleQueryString(options),
        body,
    };
    payloads.push(payload);
    return payloads;
}
exports.assembleIacLocalPayloads = assembleIacLocalPayloads;
//# sourceMappingURL=run-iac-test.js.map