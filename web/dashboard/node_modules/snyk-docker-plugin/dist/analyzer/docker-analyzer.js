"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const Debug = require("debug");
const path_1 = require("path");
const binariesAnalyzer = require("../inputs/binaries/docker");
const imageInspector = require("./image-inspector");
const osReleaseDetector = require("./os-release");
const apkInputDocker = require("../inputs/apk/docker");
const aptInputDocker = require("../inputs/apt/docker");
const rpmInputDocker = require("../inputs/rpm/docker");
const apkAnalyzer = require("./package-managers/apk");
const aptAnalyzer = require("./package-managers/apt");
const rpmAnalyzer = require("./package-managers/rpm");
const debug = Debug("snyk");
async function analyze(targetImage, dockerfileAnalysis, options) {
    try {
        await imageInspector.pullIfNotLocal(targetImage, options);
    }
    catch (error) {
        debug(`Error while running analyzer: '${error}'`);
        throw new Error("Docker error: image was not found locally and pulling failed: " +
            targetImage);
    }
    const [imageInspection, osRelease] = await Promise.all([
        imageInspector.detect(targetImage, options),
        osReleaseDetector.detectDynamically(targetImage, dockerfileAnalysis, options),
    ]);
    const [apkDbFileContent, aptDbFileContent, rpmDbFileContent,] = await Promise.all([
        apkInputDocker.getApkDbFileContent(targetImage, options),
        aptInputDocker.getAptDbFileContent(targetImage, options),
        rpmInputDocker.getRpmDbFileContent(targetImage, options),
    ]);
    let pkgManagerAnalysis;
    try {
        pkgManagerAnalysis = await Promise.all([
            apkAnalyzer.analyze(targetImage, apkDbFileContent),
            aptAnalyzer.analyze(targetImage, aptDbFileContent),
            rpmAnalyzer.analyze(targetImage, rpmDbFileContent),
        ]);
    }
    catch (error) {
        debug(`Error while running analyzer: '${error.stderr}'`);
        throw new Error("Failed to detect installed OS packages");
    }
    const { installedPackages, pkgManager } = getInstalledPackages(pkgManagerAnalysis);
    let binariesAnalysis;
    try {
        binariesAnalysis = await binariesAnalyzer.analyze(targetImage, installedPackages, pkgManager, options);
    }
    catch (error) {
        debug(`Error while running binaries analyzer: '${error}'`);
        throw new Error("Failed to detect binaries versions");
    }
    return {
        imageId: imageInspection.Id,
        osRelease,
        results: pkgManagerAnalysis,
        binaries: binariesAnalysis,
        imageLayers: imageInspection.RootFS && imageInspection.RootFS.Layers !== undefined
            ? imageInspection.RootFS.Layers.map((layer) => path_1.normalize(layer))
            : [],
    };
}
exports.analyze = analyze;
function getInstalledPackages(results) {
    const dockerAnalysis = results.find((res) => {
        return res.Analysis && res.Analysis.length > 0;
    });
    if (!dockerAnalysis) {
        return { installedPackages: [] };
    }
    const installedPackages = dockerAnalysis.Analysis.map((pkg) => pkg.Name);
    let pkgManager = dockerAnalysis.AnalyzeType;
    if (pkgManager) {
        pkgManager = pkgManager.toLowerCase();
    }
    return { installedPackages, pkgManager };
}
//# sourceMappingURL=docker-analyzer.js.map