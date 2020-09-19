"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.validateOptions = void 0;
const config = require("./config");
const reachableVulns = require("./reachable-vulns");
const is_multi_project_scan_1 = require("./is-multi-project-scan");
async function validateOptions(options, packageManager) {
    if (options.reachableVulns) {
        // Throwing error only in case when both packageManager and allProjects not defined
        if (!packageManager && !is_multi_project_scan_1.isMultiProjectScan(options)) {
            throw new Error('Could not determine package manager');
        }
        const org = options.org || config.org;
        await reachableVulns.validatePayload(org, options, packageManager);
    }
}
exports.validateOptions = validateOptions;
//# sourceMappingURL=options-validator.js.map