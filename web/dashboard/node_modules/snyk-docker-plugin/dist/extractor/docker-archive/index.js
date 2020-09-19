"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const path_1 = require("path");
var layer_1 = require("./layer");
exports.extractArchive = layer_1.extractArchive;
function getManifestLayers(manifest) {
    return manifest.Layers.map((layer) => path_1.normalize(layer));
}
exports.getManifestLayers = getManifestLayers;
function getImageIdFromManifest(manifest) {
    try {
        return manifest.Config.split(".")[0];
    }
    catch (err) {
        throw new Error("Failed to extract image ID from archive manifest");
    }
}
exports.getImageIdFromManifest = getImageIdFromManifest;
function getRootFsLayersFromConfig(imageConfig) {
    try {
        return imageConfig.rootfs.diff_ids;
    }
    catch (err) {
        throw new Error("Failed to extract rootfs array from image config");
    }
}
exports.getRootFsLayersFromConfig = getRootFsLayersFromConfig;
function getPlatformFromConfig(imageConfig) {
    return imageConfig.os && imageConfig.architecture
        ? `${imageConfig.os}/${imageConfig.architecture}`
        : undefined;
}
exports.getPlatformFromConfig = getPlatformFromConfig;
//# sourceMappingURL=index.js.map