"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const fs = require("fs");
const path = require("path");
const tmp = require("tmp");
const uuid_1 = require("uuid");
const image_inspector_1 = require("./analyzer/image-inspector");
const image_type_1 = require("./image-type");
const staticModule = require("./static");
const types_1 = require("./types");
async function experimentalAnalysis(targetImage, dockerfileAnalysis, options) {
    // assume Distroless scanning
    const imageType = image_type_1.getImageType(targetImage);
    switch (imageType) {
        case types_1.ImageType.DockerArchive:
        case types_1.ImageType.OciArchive:
            return localArchive(targetImage, imageType, dockerfileAnalysis, options);
        case types_1.ImageType.Identifier:
            return distroless(targetImage, dockerfileAnalysis, options);
        default:
            throw new Error("Unhandled image type for image " + targetImage);
    }
}
exports.experimentalAnalysis = experimentalAnalysis;
async function localArchive(targetImage, imageType, dockerfileAnalysis, options) {
    const archivePath = image_type_1.getArchivePath(targetImage);
    if (!fs.existsSync(archivePath)) {
        throw new Error("The provided archive path does not exist on the filesystem");
    }
    if (!fs.lstatSync(archivePath).isFile()) {
        throw new Error("The provided archive path is not a file");
    }
    // The target image becomes the base of the path, e.g. "archive.tar" for "/var/tmp/archive.tar"
    const imageIdentifier = path.basename(archivePath);
    return await getStaticAnalysisResult(imageIdentifier, archivePath, dockerfileAnalysis, imageType, options["app-vulns"]);
}
// experimental flow expected to be merged with the static analysis when ready
async function distroless(targetImage, dockerfileAnalysis, options) {
    var _a, _b, _c, _d;
    if (staticModule.isRequestingStaticAnalysis(options)) {
        options.staticAnalysisOptions.distroless = true;
        return staticModule.analyzeStatically(targetImage, dockerfileAnalysis, options);
    }
    const imageSavePath = fullImageSavePath((_a = options) === null || _a === void 0 ? void 0 : _a.imageSavePath);
    const archiveResult = await image_inspector_1.getImageArchive(targetImage, imageSavePath, (_b = options) === null || _b === void 0 ? void 0 : _b.username, (_c = options) === null || _c === void 0 ? void 0 : _c.password, (_d = options) === null || _d === void 0 ? void 0 : _d.platform);
    try {
        return await getStaticAnalysisResult(targetImage, archiveResult.path, dockerfileAnalysis, types_1.ImageType.DockerArchive, options["app-vulns"]);
    }
    finally {
        archiveResult.removeArchive();
    }
}
exports.distroless = distroless;
async function getStaticAnalysisResult(targetImage, archivePath, dockerfileAnalysis, imageType, appScan) {
    const scanningOptions = {
        staticAnalysisOptions: {
            imagePath: archivePath,
            imageType,
            distroless: true,
            appScan,
        },
    };
    return await staticModule.analyzeStatically(targetImage, dockerfileAnalysis, scanningOptions);
}
function fullImageSavePath(imageSavePath) {
    let imagePath = tmp.dirSync().name;
    if (imageSavePath) {
        imagePath = path.normalize(imageSavePath);
    }
    return path.join(imagePath, uuid_1.v4());
}
exports.fullImageSavePath = fullImageSavePath;
//# sourceMappingURL=experimental.js.map