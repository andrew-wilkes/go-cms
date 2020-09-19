"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const stream_1 = require("stream");
const stream_utils_1 = require("../stream-utils");
async function applyCallbacks(matchedActions, fileContentStream) {
    const result = {};
    const actionsToAwait = matchedActions.map((action) => {
        // Using a pass through allows us to read the stream multiple times.
        const streamCopy = new stream_1.PassThrough();
        fileContentStream.pipe(streamCopy);
        // Queue the promise but don't await on it yet: we want consumers to start around the same time.
        const promise = action.callback !== undefined
            ? action.callback(streamCopy)
            : // If no callback was provided for this action then return as string by default.
                stream_utils_1.streamToString(streamCopy);
        return promise.then((content) => {
            // Assign the result once the Promise is complete.
            result[action.actionName] = content;
        });
    });
    await Promise.all(actionsToAwait);
    return result;
}
exports.applyCallbacks = applyCallbacks;
//# sourceMappingURL=callbacks.js.map