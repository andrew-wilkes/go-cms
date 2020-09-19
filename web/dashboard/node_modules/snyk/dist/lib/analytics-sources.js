"use strict";
/*
  We are collecting Snyk CLI usage in our official integrations

  We distinguish them by either:
  - Setting integrationNameHeader or integrationVersionHeader in environment when CLI is run
  - passing an --integration-name or --integration-version flags on CLI invocation

  Integration name is validated with a list
*/
Object.defineProperty(exports, "__esModule", { value: true });
exports.getIntegrationVersion = exports.getIntegrationName = exports.integrationVersionHeader = exports.integrationNameHeader = void 0;
exports.integrationNameHeader = 'SNYK_INTEGRATION_NAME';
exports.integrationVersionHeader = 'SNYK_INTEGRATION_VERSION';
var TrackedIntegration;
(function (TrackedIntegration) {
    // Distribution builds/packages
    TrackedIntegration["NPM"] = "NPM";
    TrackedIntegration["STANDALONE"] = "STANDALONE";
    // tracked by passing envvar on CLI invocation
    TrackedIntegration["HOMEBREW"] = "HOMEBREW";
    TrackedIntegration["SCOOP"] = "SCOOP";
    // Our Docker images - tracked by passing envvar on CLI invocation
    TrackedIntegration["DOCKER_SNYK_CLI"] = "DOCKER_SNYK_CLI";
    TrackedIntegration["DOCKER_SNYK"] = "DOCKER_SNYK";
    // IDE plugins - tracked by passing flag or envvar on CLI invocation
    TrackedIntegration["JETBRAINS_IDE"] = "JETBRAINS_IDE";
    TrackedIntegration["ECLIPSE"] = "ECLIPSE";
    TrackedIntegration["VS_CODE_VULN_COST"] = "VS_CODE_VULN_COST";
    // CI - tracked by passing flag or envvar on CLI invocation
    TrackedIntegration["JENKINS"] = "JENKINS";
    TrackedIntegration["TEAMCITY"] = "TEAMCITY";
    TrackedIntegration["BITBUCKET_PIPELINES"] = "BITBUCKET_PIPELINES";
    TrackedIntegration["AZURE_PIPELINES"] = "AZURE_PIPELINES";
    TrackedIntegration["CIRCLECI_ORB"] = "CIRCLECI_ORB";
    TrackedIntegration["GITHUB_ACTIONS"] = "GITHUB_ACTIONS";
    // Partner integrations - tracked by passing envvar on CLI invocation
    TrackedIntegration["DOCKER_DESKTOP"] = "DOCKER_DESKTOP";
})(TrackedIntegration || (TrackedIntegration = {}));
// TODO: propagate these to the UTM params
exports.getIntegrationName = (args) => {
    var _a;
    const integrationName = String(((_a = args[0]) === null || _a === void 0 ? void 0 : _a.integrationName) || // Integration details passed through CLI flag
        process.env[exports.integrationNameHeader] ||
        '').toUpperCase();
    if (integrationName in TrackedIntegration) {
        return integrationName;
    }
    return '';
};
exports.getIntegrationVersion = (args) => {
    var _a;
    // Integration details passed through CLI flag
    const integrationVersion = String(((_a = args[0]) === null || _a === void 0 ? void 0 : _a.integrationVersion) || process.env[exports.integrationVersionHeader] || '');
    return integrationVersion;
};
//# sourceMappingURL=analytics-sources.js.map