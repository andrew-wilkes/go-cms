"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
async function getOsRelease(docker, releasePath) {
    try {
        return (await docker.catSafe(releasePath)).stdout;
    }
    catch (error) {
        throw new Error(error.stderr);
    }
}
exports.getOsRelease = getOsRelease;
//# sourceMappingURL=docker.js.map