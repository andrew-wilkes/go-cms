"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var layer_1 = require("./layer");
exports.extractArchive = layer_1.extractArchive;
function getManifestLayers(manifest) {
    return manifest.layers.map((layer) => layer.digest);
}
exports.getManifestLayers = getManifestLayers;
function getImageIdFromManifest(manifest) {
    try {
        return manifest.config.digest;
    }
    catch (err) {
        throw new Error("Failed to extract image ID from archive manifest");
    }
}
exports.getImageIdFromManifest = getImageIdFromManifest;
//# sourceMappingURL=index.js.map