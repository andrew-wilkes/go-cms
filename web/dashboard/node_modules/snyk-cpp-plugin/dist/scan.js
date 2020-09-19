"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.scan = void 0;
const fs = require("fs");
const find_1 = require("./find");
const hash_1 = require("./hash");
async function scan(options) {
    try {
        if (!options.path) {
            throw 'invalid options no path provided.';
        }
        if (!fs.existsSync(options.path)) {
            throw `'${options.path}' does not exist.`;
        }
        const filePaths = await find_1.find(options.path);
        const fingerprints = [];
        for (const filePath of filePaths) {
            const md5 = await hash_1.hash(filePath);
            fingerprints.push({
                filePath,
                hash: md5,
            });
        }
        const artifacts = [
            { type: 'cpp-fingerprints', data: fingerprints, meta: {} },
        ];
        const scanResults = [
            {
                artifacts,
                meta: {},
            },
        ];
        return scanResults;
    }
    catch (error) {
        throw new Error(`Could not scan C/C++ project, ${error}`);
    }
}
exports.scan = scan;
//# sourceMappingURL=scan.js.map