"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.pluckPolicies = void 0;
const _ = require("lodash");
function pluckPolicies(pkg) {
    if (!pkg) {
        return [];
    }
    if (pkg.snyk) {
        return pkg.snyk;
    }
    if (!pkg.dependencies) {
        return [];
    }
    return _.flatten(Object.keys(pkg.dependencies)
        .map((name) => pluckPolicies(pkg.dependencies[name]))
        .filter(Boolean));
}
exports.pluckPolicies = pluckPolicies;
//# sourceMappingURL=pluck-policies.js.map