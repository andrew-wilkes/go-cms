"use strict";
const chalk_1 = require("chalk");
const fs = require("fs");
const Debug = require("debug");
const pathUtil = require("path");
const cli_interface_1 = require("@snyk/cli-interface");
const options_validator_1 = require("../../../lib/options-validator");
const config = require("../../../lib/config");
const detect = require("../../../lib/detect");
const spinner = require("../../../lib/spinner");
const analytics = require("../../../lib/analytics");
const api_token_1 = require("../../../lib/api-token");
const print_deps_1 = require("../../../lib/print-deps");
const monitor_1 = require("../../../lib/monitor");
const process_json_monitor_1 = require("./process-json-monitor");
const snyk = require("../../../lib"); // TODO(kyegupov): fix import
const format_monitor_response_1 = require("./formatters/format-monitor-response");
const get_deps_from_plugin_1 = require("../../../lib/plugins/get-deps-from-plugin");
const get_extra_project_count_1 = require("../../../lib/plugins/get-extra-project-count");
const extract_package_manager_1 = require("../../../lib/plugins/extract-package-manager");
const convert_multi_plugin_res_to_multi_custom_1 = require("../../../lib/plugins/convert-multi-plugin-res-to-multi-custom");
const convert_single_splugin_res_to_multi_custom_1 = require("../../../lib/plugins/convert-single-splugin-res-to-multi-custom");
const dev_count_analysis_1 = require("../../../lib/monitor/dev-count-analysis");
const errors_1 = require("../../../lib/errors");
const is_multi_project_scan_1 = require("../../../lib/is-multi-project-scan");
const SEPARATOR = '\n-------------------------------------------------------\n';
const debug = Debug('snyk');
// This is used instead of `let x; try { x = await ... } catch { cleanup }` to avoid
// declaring the type of x as possibly undefined.
async function promiseOrCleanup(p, cleanup) {
    return p.catch((error) => {
        cleanup();
        throw error;
    });
}
// Returns an array of Registry responses (one per every sub-project scanned), a single response,
// or an error message.
async function monitor(...args0) {
    var _a;
    let args = [...args0];
    let monitorOptions = {};
    const results = [];
    if (typeof args[args.length - 1] === 'object') {
        monitorOptions = args.pop();
    }
    args = args.filter(Boolean);
    // populate with default path (cwd) if no path given
    if (args.length === 0) {
        args.unshift(process.cwd());
    }
    const options = monitorOptions;
    if (options.id) {
        snyk.id = options.id;
    }
    if (options.allSubProjects && options['project-name']) {
        throw new Error('`--all-sub-projects` is currently not compatible with `--project-name`');
    }
    if (options.docker && options['remote-repo-url']) {
        throw new Error('`--remote-repo-url` is not supported for container scans');
    }
    api_token_1.apiTokenExists();
    let contributors = [];
    if (!options.docker && analytics.allowAnalytics()) {
        try {
            const repoPath = process.cwd();
            const dNow = new Date();
            const timestampStartOfContributingDeveloperPeriod = dev_count_analysis_1.getTimestampStartOfContributingDevTimeframe(dNow, dev_count_analysis_1.CONTRIBUTING_DEVELOPER_PERIOD_DAYS);
            const gitLogResults = await dev_count_analysis_1.runGitLog(timestampStartOfContributingDeveloperPeriod, repoPath, dev_count_analysis_1.execShell);
            const stats = dev_count_analysis_1.parseGitLog(gitLogResults);
            contributors = stats.getRepoContributors();
        }
        catch (err) {
            debug('error getting repo contributors', err);
        }
    }
    // Part 1: every argument is a scan target; process them sequentially
    for (const path of args) {
        debug(`Processing ${path}...`);
        try {
            validateMonitorPath(path, options.docker);
            let analysisType = 'all';
            let packageManager;
            if (is_multi_project_scan_1.isMultiProjectScan(options)) {
                analysisType = 'all';
            }
            else if (options.docker) {
                analysisType = 'docker';
            }
            else {
                packageManager = detect.detectPackageManager(path, options);
            }
            await options_validator_1.validateOptions(options, packageManager);
            const targetFile = !options.scanAllUnmanaged && options.docker && !options.file // snyk monitor --docker (without --file)
                ? undefined
                : options.file || detect.detectPackageFile(path);
            const displayPath = pathUtil.relative('.', pathUtil.join(path, targetFile || ''));
            const analyzingDepsSpinnerLabel = 'Analyzing ' +
                (packageManager ? packageManager : analysisType) +
                ' dependencies for ' +
                displayPath;
            await spinner(analyzingDepsSpinnerLabel);
            // Scan the project dependencies via a plugin
            analytics.add('pluginOptions', options);
            debug('getDepsFromPlugin ...');
            // each plugin will be asked to scan once per path
            // some return single InspectResult & newer ones return Multi
            const inspectResult = await promiseOrCleanup(get_deps_from_plugin_1.getDepsFromPlugin(path, Object.assign(Object.assign({}, options), { path,
                packageManager })), spinner.clear(analyzingDepsSpinnerLabel));
            analytics.add('pluginName', inspectResult.plugin.name);
            // We send results from "all-sub-projects" scanning as different Monitor objects
            // multi result will become default, so start migrating code to always work with it
            let perProjectResult;
            if (!cli_interface_1.legacyPlugin.isMultiResult(inspectResult)) {
                perProjectResult = convert_single_splugin_res_to_multi_custom_1.convertSingleResultToMultiCustom(inspectResult);
            }
            else {
                perProjectResult = convert_multi_plugin_res_to_multi_custom_1.convertMultiResultToMultiCustom(inspectResult);
            }
            const failedResults = inspectResult
                .failedResults;
            if (failedResults === null || failedResults === void 0 ? void 0 : failedResults.length) {
                failedResults.forEach((result) => {
                    results.push({
                        ok: false,
                        data: new errors_1.MonitorError(500, result.errMessage),
                        path: result.targetFile || '',
                    });
                });
            }
            const postingMonitorSpinnerLabel = 'Posting monitor snapshot for ' + displayPath + ' ...';
            await spinner(postingMonitorSpinnerLabel);
            // Post the project dependencies to the Registry
            for (const projectDeps of perProjectResult.scannedProjects) {
                try {
                    if (!projectDeps.depGraph && !projectDeps.depTree) {
                        debug('scannedProject is missing depGraph or depTree, cannot run test/monitor');
                        throw new errors_1.FailedToRunTestError('Your monitor request could not be completed. Please email support@snyk.io');
                    }
                    const extractedPackageManager = extract_package_manager_1.extractPackageManager(projectDeps, perProjectResult, options);
                    analytics.add('packageManager', extractedPackageManager);
                    const projectName = getProjectName(projectDeps);
                    if (projectDeps.depGraph) {
                        debug(`Processing ${(_a = projectDeps.depGraph.rootPkg) === null || _a === void 0 ? void 0 : _a.name}...`);
                        print_deps_1.maybePrintDepGraph(options, projectDeps.depGraph);
                    }
                    if (projectDeps.depTree) {
                        debug(`Processing ${projectDeps.depTree.name}...`);
                        print_deps_1.maybePrintDepTree(options, projectDeps.depTree);
                    }
                    const tFile = projectDeps.targetFile || targetFile;
                    const targetFileRelativePath = projectDeps.plugin.targetFile ||
                        (tFile && pathUtil.join(pathUtil.resolve(path), tFile)) ||
                        '';
                    const res = await promiseOrCleanup(monitor_1.monitor(path, generateMonitorMeta(options, extractedPackageManager), projectDeps, options, projectDeps.plugin, targetFileRelativePath, contributors), spinner.clear(postingMonitorSpinnerLabel));
                    res.path = path;
                    const monOutput = format_monitor_response_1.formatMonitorOutput(extractedPackageManager, res, options, projectName, await get_extra_project_count_1.getExtraProjectCount(path, options, inspectResult));
                    // push a good result
                    results.push({ ok: true, data: monOutput, path, projectName });
                }
                catch (err) {
                    // pushing this error allow this inner loop to keep scanning the projects
                    // even if 1 in 100 fails
                    results.push({ ok: false, data: err, path });
                }
            }
        }
        catch (err) {
            // push this error, the loop continues
            results.push({ ok: false, data: err, path });
        }
        finally {
            spinner.clearAll();
        }
    }
    // Part 2: process the output from the Registry
    if (options.json) {
        return process_json_monitor_1.processJsonMonitorResponse(results);
    }
    const output = results
        .map((res) => {
        if (res.ok) {
            return res.data;
        }
        const errorMessage = res.data && res.data.userMessage
            ? chalk_1.default.bold.red(res.data.userMessage)
            : res.data
                ? res.data.message
                : 'Unknown error occurred.';
        return (chalk_1.default.bold.white('\nMonitoring ' + res.path + '...\n\n') + errorMessage);
    })
        .join('\n' + SEPARATOR);
    if (results.every((res) => res.ok)) {
        return output;
    }
    throw new Error(output);
}
function generateMonitorMeta(options, packageManager) {
    return {
        method: 'cli',
        packageManager,
        'policy-path': options['policy-path'],
        'project-name': options['project-name'] || config.PROJECT_NAME,
        isDocker: !!options.docker,
        prune: !!options.pruneRepeatedSubdependencies,
        'experimental-dep-graph': !!options['experimental-dep-graph'],
        'remote-repo-url': options['remote-repo-url'],
    };
}
function validateMonitorPath(path, isDocker) {
    const exists = fs.existsSync(path);
    if (!exists && !isDocker) {
        throw new Error('"' + path + '" is not a valid path for "snyk monitor"');
    }
}
function getProjectName(projectDeps) {
    var _a, _b, _c, _d;
    return (((_a = projectDeps.meta) === null || _a === void 0 ? void 0 : _a.gradleProjectName) || ((_c = (_b = projectDeps.depGraph) === null || _b === void 0 ? void 0 : _b.rootPkg) === null || _c === void 0 ? void 0 : _c.name) || ((_d = projectDeps.depTree) === null || _d === void 0 ? void 0 : _d.name));
}
module.exports = monitor;
//# sourceMappingURL=index.js.map