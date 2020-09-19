"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const Debug = require("debug");
const archiveExtractor = require("../extractor");
const static_1 = require("../inputs/apk/static");
const static_2 = require("../inputs/apt/static");
const static_3 = require("../inputs/binaries/static");
const static_4 = require("../inputs/distroless/static");
const filePatternStatic = require("../inputs/file-pattern/static");
const static_5 = require("../inputs/node/static");
const static_6 = require("../inputs/os-release/static");
const static_7 = require("../inputs/rpm/static");
const types_1 = require("../types");
const applications_1 = require("./applications");
const osReleaseDetector = require("./os-release");
const apk_1 = require("./package-managers/apk");
const apt_1 = require("./package-managers/apt");
const rpm_1 = require("./package-managers/rpm");
const debug = Debug("snyk");
const supportedArchives = [
    types_1.ImageType.DockerArchive,
    types_1.ImageType.OciArchive,
];
async function analyze(targetImage, dockerfileAnalysis, options) {
    if (!supportedArchives.includes(options.imageType)) {
        throw new Error("Unhandled image type");
    }
    const staticAnalysisActions = [
        static_1.getApkDbFileContentAction,
        static_2.getDpkgFileContentAction,
        static_2.getExtFileContentAction,
        static_7.getRpmDbFileContentAction,
        ...static_6.getOsReleaseActions,
        static_3.getNodeBinariesFileContentAction,
        static_3.getOpenJDKBinariesFileContentAction,
        static_5.getNodeAppFileContentAction,
    ];
    if (options.distroless) {
        staticAnalysisActions.push(static_4.getDpkgPackageFileContentAction);
    }
    const checkForGlobs = shouldCheckForGlobs(options);
    if (checkForGlobs) {
        staticAnalysisActions.push(filePatternStatic.generateExtractAction(options.globsToFind.include, options.globsToFind.exclude));
    }
    const { imageId, manifestLayers, extractedLayers, rootFsLayers, platform, } = await archiveExtractor.extractImageContent(options.imageType, options.imagePath, staticAnalysisActions);
    const [apkDbFileContent, aptDbFileContent, rpmDbFileContent,] = await Promise.all([
        static_1.getApkDbFileContent(extractedLayers),
        static_2.getAptDbFileContent(extractedLayers),
        static_7.getRpmDbFileContent(extractedLayers),
    ]);
    let distrolessAptFiles = [];
    if (options.distroless) {
        distrolessAptFiles = static_4.getAptFiles(extractedLayers);
    }
    const manifestFiles = [];
    if (checkForGlobs) {
        const matchingFiles = filePatternStatic.getMatchingFiles(extractedLayers);
        manifestFiles.push(...matchingFiles);
    }
    let osRelease;
    try {
        osRelease = await osReleaseDetector.detectStatically(extractedLayers, dockerfileAnalysis);
    }
    catch (err) {
        debug(err);
        throw new Error("Failed to detect OS release");
    }
    let results;
    try {
        results = await Promise.all([
            apk_1.analyze(targetImage, apkDbFileContent),
            apt_1.analyze(targetImage, aptDbFileContent),
            rpm_1.analyze(targetImage, rpmDbFileContent),
            apt_1.analyzeDistroless(targetImage, distrolessAptFiles),
        ]);
    }
    catch (err) {
        debug(err);
        throw new Error("Failed to detect installed OS packages");
    }
    const binaries = static_3.getBinariesHashes(extractedLayers);
    const applicationDependenciesScanResults = [];
    if (options.appScan) {
        const nodeDependenciesScanResults = await applications_1.nodeFilesToScannedProjects(static_5.getNodeAppFileContent(extractedLayers));
        applicationDependenciesScanResults.push(...nodeDependenciesScanResults);
    }
    return {
        imageId,
        osRelease,
        platform,
        results,
        binaries,
        imageLayers: manifestLayers,
        rootFsLayers,
        applicationDependenciesScanResults,
        manifestFiles,
    };
}
exports.analyze = analyze;
function shouldCheckForGlobs(options) {
    return (options &&
        options.globsToFind &&
        options.globsToFind.include &&
        Array.isArray(options.globsToFind.include) &&
        options.globsToFind.include.length > 0);
}
//# sourceMappingURL=static-analyzer.js.map