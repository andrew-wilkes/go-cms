"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const semver_1 = require("semver");
const docker_1 = require("../../../docker");
async function extract(targetImage, options) {
    try {
        const binaryVersion = (await new docker_1.Docker(targetImage, options).runSafe("node", ["--version"])).stdout;
        return parseNodeBinary(binaryVersion);
    }
    catch (error) {
        throw new Error(error.stderr);
    }
}
exports.extract = extract;
function parseNodeBinary(version) {
    const nodeVersion = semver_1.valid(version && version.trim());
    if (!nodeVersion) {
        return null;
    }
    return {
        name: "node",
        version: nodeVersion,
    };
}
const packageNames = ["node", "nodejs"];
function installedByPackageManager(installedPackages) {
    return (installedPackages.filter((pkg) => packageNames.indexOf(pkg) > -1).length > 0);
}
exports.installedByPackageManager = installedByPackageManager;
//# sourceMappingURL=node.js.map