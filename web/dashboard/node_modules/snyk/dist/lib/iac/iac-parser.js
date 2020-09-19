"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateK8sFile = void 0;
//TODO(orka): take out into a new lib
const YAML = require("js-yaml");
const debugLib = require("debug");
const errors_1 = require("../errors");
const debug = debugLib('snyk-detect');
const mandatoryKeysForSupportedK8sKinds = {
    deployment: ['apiVersion', 'metadata', 'spec'],
    pod: ['apiVersion', 'metadata', 'spec'],
    service: ['apiVersion', 'metadata', 'spec'],
    podsecuritypolicy: ['apiVersion', 'metadata', 'spec'],
    networkpolicy: ['apiVersion', 'metadata', 'spec'],
};
function getFileType(filePath) {
    const filePathSplit = filePath.split('.');
    return filePathSplit[filePathSplit.length - 1].toLowerCase();
}
function parseYamlOrJson(fileContent, filePath) {
    const fileType = getFileType(filePath);
    switch (fileType) {
        case 'yaml':
        case 'yml':
            try {
                return YAML.safeLoadAll(fileContent);
            }
            catch (e) {
                debug('Failed to parse iac config as a YAML');
            }
            break;
        case 'json':
            try {
                const objectsArr = [];
                objectsArr.push(JSON.parse(fileContent));
                return objectsArr;
            }
            catch (e) {
                debug('Failed to parse iac config as a JSON');
            }
            break;
        default:
            debug(`Unsupported iac config file type (${fileType})`);
    }
    return undefined;
}
// This function validates that there is at least one valid doc with a k8s object kind.
// A valid k8s object has a kind key (.kind) from the keys of `mandatoryKeysForSupportedK8sKinds`
// and all of the keys from `mandatoryKeysForSupportedK8sKinds[kind]`.
// If there is a doc with a supported kind, but invalid, we should fail
// The function return true if the yaml is a valid k8s one, or false otherwise
function validateK8sFile(fileContent, filePath, root) {
    const k8sObjects = parseYamlOrJson(fileContent, filePath);
    if (!k8sObjects) {
        throw errors_1.IllegalIacFileError([root]);
    }
    let numOfSupportedKeyDocs = 0;
    for (let i = 0; i < k8sObjects.length; i++) {
        const k8sObject = k8sObjects[i];
        if (!k8sObject || !k8sObject.kind) {
            continue;
        }
        const kind = k8sObject.kind.toLowerCase();
        if (!Object.keys(mandatoryKeysForSupportedK8sKinds).includes(kind)) {
            continue;
        }
        numOfSupportedKeyDocs++;
        for (let i = 0; i < mandatoryKeysForSupportedK8sKinds[kind].length; i++) {
            const key = mandatoryKeysForSupportedK8sKinds[kind][i];
            if (!k8sObject[key]) {
                debug(`Missing key (${key}) from supported k8s object kind (${kind})`);
                throw errors_1.IllegalIacFileError([root]);
            }
        }
    }
    if (numOfSupportedKeyDocs === 0) {
        throw errors_1.NotSupportedIacFileError([root]);
    }
    debug(`k8s config found (${filePath})`);
}
exports.validateK8sFile = validateK8sFile;
//# sourceMappingURL=iac-parser.js.map