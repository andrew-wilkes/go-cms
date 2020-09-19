"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const gunzip = require("gunzip-maybe");
const path = require("path");
const tar_stream_1 = require("tar-stream");
const callbacks_1 = require("./callbacks");
/**
 * Extract key files from the specified TAR stream.
 * @param layerTarStream image layer as a Readable TAR stream. Note: consumes the stream.
 * @param extractActions array of pattern, callbacks pairs
 * @returns extracted file products
 */
async function extractImageLayer(layerTarStream, extractActions) {
    return new Promise((resolve, reject) => {
        const result = {};
        const tarExtractor = tar_stream_1.extract();
        tarExtractor.on("entry", async (headers, stream, next) => {
            if (headers.type === "file") {
                const absoluteFileName = path.join(path.sep, headers.name);
                // TODO wouldn't it be simpler to first check
                // if the filename matches any patterns?
                const processedResult = await extractFileAndProcess(absoluteFileName, stream, extractActions);
                if (processedResult !== undefined) {
                    result[absoluteFileName] = processedResult;
                }
            }
            stream.resume(); // auto drain the stream
            next(); // ready for next entry
        });
        tarExtractor.on("finish", () => {
            // all layer level entries read
            resolve(result);
        });
        tarExtractor.on("error", (error) => reject(error));
        layerTarStream.pipe(gunzip()).pipe(tarExtractor);
    });
}
exports.extractImageLayer = extractImageLayer;
/**
 * Note: consumes the stream.
 */
async function extractFileAndProcess(fileName, fileStream, extractActions) {
    const matchedActions = extractActions.filter((action) => action.filePathMatches(fileName));
    if (matchedActions.length > 0) {
        return await callbacks_1.applyCallbacks(matchedActions, fileStream);
    }
    return undefined;
}
//# sourceMappingURL=layer.js.map