"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getManifest = void 0;
const registry_call_1 = require("./registry-call");
const contentTypes = require("./content-types");
const needle_1 = require("./needle");
async function getManifest(registryBase, repo, tag, username, password, options = {}) {
    const endpoint = `/${repo}/manifests/${tag}`;
    const manifestResponse = await registry_call_1.registryV2Call(registryBase, endpoint, contentTypes.MANIFEST_V2, username, password, options);
    return needle_1.parseResponseBody(manifestResponse);
}
exports.getManifest = getManifest;
//# sourceMappingURL=get-manifest.js.map