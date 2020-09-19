"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.apiTokenExists = exports.getDockerToken = exports.api = void 0;
const errors_1 = require("../lib/errors");
const config = require("./config");
const user_config_1 = require("./user-config");
function api() {
    // note: config.TOKEN will potentially come via the environment
    return config.api || config.TOKEN || user_config_1.config.get('api');
}
exports.api = api;
function getDockerToken() {
    return process.env.SNYK_DOCKER_TOKEN;
}
exports.getDockerToken = getDockerToken;
function apiTokenExists() {
    const configured = api();
    if (!configured) {
        throw new errors_1.MissingApiTokenError();
    }
    return configured;
}
exports.apiTokenExists = apiTokenExists;
//# sourceMappingURL=api-token.js.map