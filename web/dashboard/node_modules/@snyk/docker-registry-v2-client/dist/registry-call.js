"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getChallengeURL = exports.paginatedV2Call = exports.registryV2Call = void 0;
const parseLink = require("parse-link-header");
const path = require("path");
const url = require("url");
const needle_1 = require("./needle");
const BEARER_REALM = "Bearer realm";
async function registryV2Call(registryBase, endpoint, accept, username, password, reqOptions = {}) {
    const reqConfig = applyRequestOptions({
        headers: { Accept: accept },
        uri: `https://${registryBase}/v2${endpoint}`,
    }, reqOptions);
    try {
        return await needle_1.needleWrapper(reqConfig);
    }
    catch (err) {
        if (err.statusCode === 401) {
            if (!username || !password) {
                // Supply and empty username and password if no credentials
                // are provided. These might be added later by a broker client.
                username = username ? username : "";
                password = password ? password : "";
            }
            const authConfig = await setAuthConfig(registryBase, err, reqConfig, username, password, reqOptions);
            try {
                return await needle_1.needleWrapper(authConfig);
            }
            catch (err) {
                if (err.statusCode === 307 || err.statusCode === 302) {
                    return await handleRedirect(err, reqConfig);
                }
                throw err;
            }
        }
        if (err.statusCode === 307 || err.statusCode === 302) {
            return await handleRedirect(err, reqConfig);
        }
        throw err;
    }
}
exports.registryV2Call = registryV2Call;
async function paginatedV2Call(registryBase, accept, username, password, endpoint, key, pageSize = 1000, maxPages = Number.MAX_SAFE_INTEGER, reqOptions = {}) {
    const result = [];
    let pageEndpoint = `${endpoint}?n=${pageSize}`;
    let pageCount = 0;
    while (pageCount < maxPages) {
        const response = await registryV2Call(registryBase, pageEndpoint, accept, username, password, reqOptions);
        const body = needle_1.parseResponseBody(response);
        result.push(...body[key]);
        if (!response.headers.link) {
            break;
        }
        pageCount += 1;
        pageEndpoint = pageEndpointForLink(endpoint, response.headers.link);
    }
    return result;
}
exports.paginatedV2Call = paginatedV2Call;
function pageEndpointForLink(endpoint, link) {
    const linkPath = parseLink(link).next.url;
    const linkQuery = linkPath.split("?")[1];
    return `${endpoint}?${linkQuery}`;
}
// Construct the challenge URL. If there is a base URL path the challenge must be relative to it.
// e.g. if the base URL is broker-container.snyk.io/broker/123 and the server returns the challenge
// URL as broker-container.snyk.io/v2/token then the correct challenge URL is:
// broker-container.snyk.io/broker/123/v2/token
function getChallengeURL(challenge, registryBase) {
    const registryBaseURL = url.parse(registryBase);
    // this assumes we are handling brokered request
    if (registryBaseURL.path !== "/") {
        const challengeURL = url.parse(challenge);
        challengeURL.pathname = challengeURL.path = path.posix.join(registryBaseURL.path, challengeURL.path);
        challenge = url.format(challengeURL);
    }
    return challenge;
}
exports.getChallengeURL = getChallengeURL;
async function getToken(registryBase, authBase, service, scope, username, password, reqOptions = {}) {
    const reqConfig = applyRequestOptions({
        uri: getChallengeURL(authBase, `https://${registryBase}/`),
        qs: {
            service,
            scope,
        },
    }, Object.assign({}, reqOptions));
    // Test truthiness, should be false when username and password are undefined
    if (username && password) {
        reqConfig.username = username;
        reqConfig.password = password;
    }
    const response = await needle_1.needleWrapper(reqConfig);
    const body = needle_1.parseResponseBody(response);
    return body.token || body.access_token;
}
async function setAuthConfig(registryBase, err, reqConfig, username, password, reqOptions) {
    // See: https://docs.docker.com/registry/spec/auth/token/#how-to-authenticate
    const challengeHeaders = err.headers["www-authenticate"];
    if (!challengeHeaders) {
        throw err;
    }
    const [authBase, service, scope] = parseChallengeHeaders(challengeHeaders);
    if (!authBase) {
        // basic auth
        return Object.assign(Object.assign({}, reqConfig), { username, password });
    }
    else {
        // bearer token
        const token = await getToken(registryBase, authBase, service, scope, username, password, reqOptions);
        return Object.assign(Object.assign({}, reqConfig), { headers: Object.assign(Object.assign({}, reqConfig.headers), { Authorization: `Bearer ${token}` }) });
    }
}
function parseChallengeHeaders(challangeHeaders) {
    const headersMap = challangeHeaders.split(",").reduce((map, entry) => {
        const [key, value] = entry.split("=");
        map[key] = JSON.parse(value);
        return map;
    }, {});
    return [headersMap[BEARER_REALM], headersMap.service, headersMap.scope];
}
async function handleRedirect(err, config) {
    // ACR does not handle redirects well, where automatic redirects
    // fail due to an unexpected authorization header.
    // the solution is to follow the redirect, however discarding
    // the token.
    const location = err.headers.location;
    if (!location) {
        throw err;
    }
    // Only clear the Authorization headers if the redirect is for
    // azure container registries.
    if (location.includes("azurecr.io")) {
        config.headers.Authorization = undefined;
    }
    config.uri = location;
    return await needle_1.needleWrapper(config);
}
/*
 * Takes request config and applies allowed options to it.
 * @param reqConfig - request config that is passed to the request library.
 * @param reqOptions - options passed in from outside of v2 client library.
 */
function applyRequestOptions(reqConfig, reqOptions) {
    const options = Object.assign({}, reqOptions);
    let uri = applyUriProtocol(reqConfig.uri, options.protocol);
    delete options.protocol;
    uri = applyUriHostMappings(uri, options.hostMappings);
    delete options.hostMappings;
    const headers = applyHeaders(reqConfig.headers, options.headers);
    delete options.headers;
    return Object.assign(Object.assign(Object.assign({}, reqConfig), options), { uri,
        headers });
}
function applyUriProtocol(uri, protocol) {
    if (!protocol) {
        return uri;
    }
    const updatedUrl = url.parse(uri);
    updatedUrl.protocol = protocol;
    return url.format(updatedUrl);
}
/**
 * Applies host mappings to given uri.
 *
 * @param uri
 * @param mappings - Array of mappings. Each mapping is represented as array
 *                   tuple: [host_regex_matcher, new_host].
 */
function applyUriHostMappings(uri, mappings) {
    if (!mappings) {
        return uri;
    }
    const updatedUrl = url.parse(uri);
    const mapping = mappings.find(([matcher]) => updatedUrl.host.match(matcher));
    if (!mapping) {
        return uri;
    }
    updatedUrl.host = mapping[1];
    return url.format(updatedUrl);
}
function applyHeaders(currentHeaders, addHeaders) {
    return Object.assign(Object.assign({}, (currentHeaders || {})), (addHeaders || {}));
}
//# sourceMappingURL=registry-call.js.map