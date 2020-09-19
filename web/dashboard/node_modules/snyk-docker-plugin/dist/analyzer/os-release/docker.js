"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const docker_1 = require("../../docker");
const os_release_1 = require("../../inputs/os-release");
const types_1 = require("../../types");
const release_analyzer_1 = require("./release-analyzer");
async function detect(targetImage, dockerfileAnalysis, options) {
    const docker = new docker_1.Docker(targetImage, options);
    let osRelease = await os_release_1.getOsReleaseDynamic(docker, types_1.OsReleaseFilePath.Linux).then((release) => release_analyzer_1.tryOSRelease(release));
    // First generic fallback
    if (!osRelease) {
        osRelease = await os_release_1.getOsReleaseDynamic(docker, types_1.OsReleaseFilePath.Lsb).then((release) => release_analyzer_1.tryLsbRelease(release));
    }
    // Fallbacks for specific older distributions
    if (!osRelease) {
        osRelease = await os_release_1.getOsReleaseDynamic(docker, types_1.OsReleaseFilePath.Debian).then((release) => release_analyzer_1.tryDebianVersion(release));
    }
    if (!osRelease) {
        osRelease = await os_release_1.getOsReleaseDynamic(docker, types_1.OsReleaseFilePath.Alpine).then((release) => release_analyzer_1.tryAlpineRelease(release));
    }
    if (!osRelease) {
        osRelease = await os_release_1.getOsReleaseDynamic(docker, types_1.OsReleaseFilePath.Oracle).then((release) => release_analyzer_1.tryOracleRelease(release));
    }
    if (!osRelease) {
        osRelease = await os_release_1.getOsReleaseDynamic(docker, types_1.OsReleaseFilePath.RedHat).then((release) => release_analyzer_1.tryRedHatRelease(release));
    }
    if (!osRelease) {
        if (dockerfileAnalysis && dockerfileAnalysis.baseImage === "scratch") {
            // If the docker file was build from a scratch image
            // then we don't have a known OS
            osRelease = { name: "scratch", version: "0.0", prettyName: "" };
        }
        else {
            osRelease = { name: "unknown", version: "0.0", prettyName: "" };
        }
    }
    // Oracle Linux identifies itself as "ol"
    if (osRelease.name.trim() === "ol") {
        osRelease.name = "oracle";
    }
    return osRelease;
}
exports.detect = detect;
//# sourceMappingURL=docker.js.map