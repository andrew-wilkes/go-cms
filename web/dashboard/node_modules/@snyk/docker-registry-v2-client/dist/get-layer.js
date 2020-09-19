"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getLayer = void 0;
const registry_call_1 = require("./registry-call");
const contentTypes = require("./content-types");
async function getLayer(registryBase, repo, digest, username, password, options = {}) {
    const endpoint = `/${repo}/blobs/${digest}`;
    options = Object.assign({ json: false, encoding: null }, options);
    const layerResponse = await registry_call_1.registryV2Call(registryBase, endpoint, contentTypes.LAYER, username, password, options);
    return layerResponse.body;
}
exports.getLayer = getLayer;
//# sourceMappingURL=get-layer.js.map