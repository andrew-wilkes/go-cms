"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const os_release_1 = require("../../inputs/os-release");
const types_1 = require("../../types");
const release_analyzer_1 = require("./release-analyzer");
async function detect(extractedLayers, dockerfileAnalysis) {
    let osRelease = await release_analyzer_1.tryOSRelease(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.Linux));
    // Fallback for the case where the same file exists in different location
    // or is a symlink to the other location
    if (!osRelease) {
        osRelease = await release_analyzer_1.tryOSRelease(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.LinuxFallback));
    }
    // Generic fallback
    if (!osRelease) {
        osRelease = await release_analyzer_1.tryLsbRelease(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.Lsb));
    }
    // Fallbacks for specific older distributions
    if (!osRelease) {
        osRelease = await release_analyzer_1.tryDebianVersion(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.Debian));
    }
    if (!osRelease) {
        osRelease = await release_analyzer_1.tryAlpineRelease(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.Alpine));
    }
    if (!osRelease) {
        osRelease = await release_analyzer_1.tryOracleRelease(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.Oracle));
    }
    if (!osRelease) {
        osRelease = await release_analyzer_1.tryRedHatRelease(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.RedHat));
    }
    if (!osRelease) {
        osRelease = await release_analyzer_1.tryCentosRelease(os_release_1.getOsReleaseStatic(extractedLayers, types_1.OsReleaseFilePath.Centos));
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
//# sourceMappingURL=static.js.map