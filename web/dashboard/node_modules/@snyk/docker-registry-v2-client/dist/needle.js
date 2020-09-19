"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.parseResponseBody = exports.needleWrapper = void 0;
const needle = require("needle");
/**
 * A wrapper function that uses `needle` for making HTTP requests,
 * and returns a response that matches what the response it used to get from `request` library
 * @param options request options
 */
async function needleWrapper(options) {
    var _a, _b;
    let uri = options.uri;
    // append query parameters
    if (options.qs) {
        for (const key in options.qs) {
            if (options.qs[key] !== undefined) {
                uri += `&${key}=${options.qs[key]}`;
            }
        }
        if (!uri.includes("?")) {
            uri = uri.replace("&", "?");
        }
    }
    const response = await needle("get", uri, options);
    // throw an error in case status code is not 2xx
    if (response && response.statusCode >= 300) {
        let message;
        if (((_b = (_a = response.body) === null || _a === void 0 ? void 0 : _a.errors) === null || _b === void 0 ? void 0 : _b.length) > 0) {
            message = response.body.errors[0].message;
        }
        else {
            message = response.body;
        }
        throw new NeedleWrapperException(message, response.statusCode, response.headers);
    }
    return response;
}
exports.needleWrapper = needleWrapper;
function parseResponseBody(response) {
    let body;
    try {
        body = JSON.parse(response.body);
    }
    catch (err) {
        body = response.body;
    }
    return body;
}
exports.parseResponseBody = parseResponseBody;
class NeedleWrapperException extends Error {
    constructor(message, statusCode, headers) {
        super(message);
        this.statusCode = statusCode;
        this.headers = headers;
    }
}
//# sourceMappingURL=needle.js.map