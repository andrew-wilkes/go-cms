"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.mapIacTestResponseToSarifResults = exports.mapIacTestResponseToSarifTool = exports.createSarifOutputForIac = exports.capitalizePackageManager = exports.getIacDisplayedOutput = void 0;
const chalk_1 = require("chalk");
const Debug = require("debug");
const formatters_1 = require("./formatters");
const remediation_based_format_issues_1 = require("./formatters/remediation-based-format-issues");
const legacy_format_issue_1 = require("./formatters/legacy-format-issue");
const legacy_1 = require("../../../lib/snyk-test/legacy");
const upperFirst = require("lodash/upperFirst");
const debug = Debug('iac-output');
function formatIacIssue(issue, isNew, path) {
    const severitiesColourMapping = {
        low: {
            colorFunc(text) {
                return chalk_1.default.blueBright(text);
            },
        },
        medium: {
            colorFunc(text) {
                return chalk_1.default.yellowBright(text);
            },
        },
        high: {
            colorFunc(text) {
                return chalk_1.default.redBright(text);
            },
        },
    };
    const newBadge = isNew ? ' (new)' : '';
    const name = issue.subType ? ` in ${chalk_1.default.bold(issue.subType)}` : '';
    let introducedBy = '';
    if (path) {
        // In this mode, we show only one path by default, for compactness
        const pathStr = remediation_based_format_issues_1.printPath(path);
        introducedBy = `\n    introduced by ${pathStr}`;
    }
    const description = extractOverview(issue.description).trim();
    const descriptionLine = `\n    ${description}\n`;
    return (severitiesColourMapping[issue.severity].colorFunc(`  ✗ ${chalk_1.default.bold(issue.title)}${newBadge} [${legacy_format_issue_1.titleCaseText(issue.severity)} Severity]`) +
        ` [${issue.id}]` +
        name +
        introducedBy +
        descriptionLine);
}
function extractOverview(description) {
    if (!description) {
        return '';
    }
    const overviewRegExp = /## Overview([\s\S]*?)(?=##|(# Details))/m;
    const overviewMatches = overviewRegExp.exec(description);
    return (overviewMatches && overviewMatches[1]) || '';
}
function getIacDisplayedOutput(iacTest, testedInfoText, meta, prefix) {
    const issuesTextArray = [
        chalk_1.default.bold.white('\nInfrastructure as code issues:'),
    ];
    const NotNew = false;
    const issues = iacTest.result.cloudConfigResults;
    debug(`iac display output - ${issues.length} issues`);
    issues
        .sort((a, b) => formatters_1.getSeverityValue(b.severity) - formatters_1.getSeverityValue(a.severity))
        .forEach((issue) => {
        issuesTextArray.push(formatIacIssue(issue, NotNew, issue.cloudConfigPath));
    });
    const issuesInfoOutput = [];
    debug(`Iac display output - ${issuesTextArray.length} issues text`);
    if (issuesTextArray.length > 0) {
        issuesInfoOutput.push(issuesTextArray.join('\n'));
    }
    let body = issuesInfoOutput.join('\n\n') + '\n\n' + meta;
    const vulnCountText = `found ${issues.length} issues`;
    const summary = testedInfoText + ', ' + chalk_1.default.red.bold(vulnCountText);
    body = body + '\n\n' + summary;
    return prefix + body;
}
exports.getIacDisplayedOutput = getIacDisplayedOutput;
function capitalizePackageManager(type) {
    switch (type) {
        case 'k8sconfig': {
            return 'Kubernetes';
        }
        case 'helmconfig': {
            return 'Helm';
        }
        case 'terraformconfig': {
            return 'Terraform';
        }
        default: {
            return 'Infrastracture as Code';
        }
    }
}
exports.capitalizePackageManager = capitalizePackageManager;
function createSarifOutputForIac(iacTestResponses) {
    const sarifRes = {
        version: '2.1.0',
        runs: [],
    };
    iacTestResponses
        .filter((iacTestResponse) => { var _a; return (_a = iacTestResponse.result) === null || _a === void 0 ? void 0 : _a.cloudConfigResults; })
        .forEach((iacTestResponse) => {
        sarifRes.runs.push({
            tool: mapIacTestResponseToSarifTool(iacTestResponse),
            results: mapIacTestResponseToSarifResults(iacTestResponse),
        });
    });
    return sarifRes;
}
exports.createSarifOutputForIac = createSarifOutputForIac;
function getIssueLevel(severity) {
    return severity === legacy_1.SEVERITY.HIGH ? 'error' : 'warning';
}
function mapIacTestResponseToSarifTool(iacTestResponse) {
    const tool = {
        driver: {
            name: 'Snyk Infrastructure as Code',
            rules: [],
        },
    };
    const pushedIds = {};
    iacTestResponse.result.cloudConfigResults.forEach((iacIssue) => {
        var _a;
        if (pushedIds[iacIssue.id]) {
            return;
        }
        (_a = tool.driver.rules) === null || _a === void 0 ? void 0 : _a.push({
            id: iacIssue.id,
            shortDescription: {
                text: `${upperFirst(iacIssue.severity)} - ${iacIssue.title}`,
            },
            fullDescription: {
                text: `Kubernetes ${iacIssue.subType}`,
            },
            help: {
                text: '',
                markdown: iacIssue.description,
            },
            defaultConfiguration: {
                level: getIssueLevel(iacIssue.severity),
            },
            properties: {
                tags: ['security', `kubernetes/${iacIssue.subType}`],
            },
        });
        pushedIds[iacIssue.id] = true;
    });
    return tool;
}
exports.mapIacTestResponseToSarifTool = mapIacTestResponseToSarifTool;
function mapIacTestResponseToSarifResults(iacTestResponse) {
    return iacTestResponse.result.cloudConfigResults.map((iacIssue) => ({
        ruleId: iacIssue.id,
        message: {
            text: `This line contains a potential ${iacIssue.severity} severity misconfiguration affacting the Kubernetes ${iacIssue.subType}`,
        },
        locations: [
            {
                physicalLocation: {
                    artifactLocation: {
                        uri: iacTestResponse.targetFile,
                    },
                    region: {
                        startLine: iacIssue.lineNumber,
                    },
                },
            },
        ],
    }));
}
exports.mapIacTestResponseToSarifResults = mapIacTestResponseToSarifResults;
//# sourceMappingURL=iac-output.js.map