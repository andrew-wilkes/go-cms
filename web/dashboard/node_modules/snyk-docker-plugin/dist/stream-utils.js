"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const crypto = require("crypto");
const HASH_ALGORITHM = "sha256"; // TODO algorithm?
const HASH_ENCODING = "hex";
async function streamToString(stream, encoding = "utf8") {
    const chunks = [];
    return new Promise((resolve, reject) => {
        stream.on("end", () => {
            resolve(chunks.join(""));
        });
        stream.on("error", (error) => reject(error));
        stream.on("data", (chunk) => {
            chunks.push(chunk.toString(encoding));
        });
    });
}
exports.streamToString = streamToString;
async function streamToBuffer(stream) {
    const chunks = [];
    return new Promise((resolve, reject) => {
        stream.on("end", () => {
            resolve(Buffer.concat(chunks));
        });
        stream.on("error", (error) => reject(error));
        stream.on("data", (chunk) => {
            chunks.push(Buffer.from(chunk));
        });
    });
}
exports.streamToBuffer = streamToBuffer;
async function streamToHash(stream) {
    return new Promise((resolve, reject) => {
        const hash = crypto.createHash(HASH_ALGORITHM);
        hash.setEncoding(HASH_ENCODING);
        stream.on("end", () => {
            hash.end();
            resolve(hash.read());
        });
        stream.on("error", (error) => reject(error));
        stream.pipe(hash);
    });
}
exports.streamToHash = streamToHash;
async function streamToJson(stream) {
    const file = await streamToString(stream);
    return JSON.parse(file);
}
exports.streamToJson = streamToJson;
//# sourceMappingURL=stream-utils.js.map