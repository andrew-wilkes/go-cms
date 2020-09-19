"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const types_1 = require("../../../analyzer/types");
async function analyze(targetImage, installedPackages, pkgManager, options) {
    const binaries = await getBinaries(targetImage, installedPackages, pkgManager, options);
    return {
        Image: targetImage,
        AnalyzeType: types_1.AnalysisType.Binaries,
        Analysis: binaries,
    };
}
exports.analyze = analyze;
const binaryVersionExtractors = {
    node: require("./node"),
    openjdk: require("./openjdk-jre"),
};
async function getBinaries(targetImage, installedPackages, pkgManager, options) {
    const binaries = [];
    for (const versionExtractor of Object.keys(binaryVersionExtractors)) {
        const extractor = binaryVersionExtractors[versionExtractor];
        if (extractor.installedByPackageManager(installedPackages, pkgManager, options)) {
            continue;
        }
        const binary = await extractor.extract(targetImage, options);
        if (binary) {
            binaries.push(binary);
        }
    }
    return binaries;
}
//# sourceMappingURL=index.js.map