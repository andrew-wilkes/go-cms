"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.find = void 0;
const fs = require("fs");
const path = require("path");
const types_1 = require("./types");
const debug_1 = require("debug");
const debug = debug_1.default('snyk-cpp-plugin');
async function find(dir) {
    const result = [];
    const dirStat = fs.statSync(dir);
    if (!dirStat.isDirectory()) {
        throw new Error(`${dir} is not a directory`);
    }
    const paths = fs.readdirSync(dir);
    for (const relativePath of paths) {
        const absolutePath = path.join(dir, relativePath);
        try {
            const stat = fs.statSync(absolutePath);
            if (stat.isDirectory()) {
                const subFiles = await find(absolutePath);
                result.push(...subFiles);
            }
            if (stat.isFile()) {
                const ext = path.extname(absolutePath);
                if (types_1.SupportFileExtensions.includes(ext)) {
                    result.push(absolutePath);
                }
            }
        }
        catch (error) {
            debug(error.message || `Error reading file ${relativePath}. ${error}`);
        }
    }
    return result;
}
exports.find = find;
//# sourceMappingURL=find.js.map