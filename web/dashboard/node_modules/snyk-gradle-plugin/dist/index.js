"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.exportsForTests = exports.buildGraph = exports.inspect = void 0;
const os = require("os");
const fs = require("fs");
const path = require("path");
const subProcess = require("./sub-process");
const tmp = require("tmp");
const errors_1 = require("./errors");
const chalk = require("chalk");
const dep_graph_1 = require("@snyk/dep-graph");
const cli_interface_1 = require("@snyk/cli-interface");
const debugModule = require("debug");
const find_cycles_1 = require("./find-cycles");
// To enable debugging output, use `snyk -d`
let logger = null;
function debugLog(s) {
    if (logger === null) {
        // Lazy init: Snyk CLI needs to process the CLI argument "-d" first.
        // TODO(BST-648): more robust handling of the debug settings
        if (process.env.DEBUG) {
            debugModule.enable(process.env.DEBUG);
        }
        logger = debugModule('snyk-gradle-plugin');
    }
    logger(s);
}
const isWin = /^win/.test(os.platform());
const quot = isWin ? '"' : "'";
const cannotResolveVariantMarkers = [
    'Cannot choose between the following',
    'Could not select value from candidates',
    'Unable to find a matching variant of project',
];
// General implementation. The result type depends on the runtime type of `options`.
async function inspect(root, targetFile, options) {
    var _a, _b;
    debugLog('Gradle inspect called with: ' +
        JSON.stringify({
            root,
            targetFile,
            allSubProjects: (_a = options) === null || _a === void 0 ? void 0 : _a.allSubProjects,
            subProject: (_b = options) === null || _b === void 0 ? void 0 : _b.subProject,
        }));
    if (!options) {
        options = { dev: false };
    }
    let subProject = options.subProject;
    if (subProject) {
        subProject = subProject.trim();
    }
    const plugin = {
        name: 'bundled:gradle',
        runtime: 'unknown',
        targetFile: targetFileFilteredForCompatibility(targetFile),
        meta: {},
    };
    if (cli_interface_1.legacyPlugin.isMultiSubProject(options)) {
        if (subProject) {
            throw new Error('gradle-sub-project flag is incompatible with multiDepRoots');
        }
        const scannedProjects = await getAllDepsAllProjects(root, targetFile, options);
        plugin.meta = plugin.meta || {};
        return {
            plugin,
            scannedProjects,
        };
    }
    const depGraphAndDepRootNames = await getAllDepsOneProject(root, targetFile, options, subProject);
    if (depGraphAndDepRootNames.allSubProjectNames) {
        plugin.meta = plugin.meta || {};
        plugin.meta.allSubProjectNames = depGraphAndDepRootNames.allSubProjectNames;
    }
    return {
        plugin,
        package: null,
        dependencyGraph: depGraphAndDepRootNames.depGraph,
        meta: {
            gradleProjectName: depGraphAndDepRootNames.gradleProjectName,
            versionBuildInfo: depGraphAndDepRootNames.versionBuildInfo,
        },
    };
}
exports.inspect = inspect;
// See the comment for DepRoot.targetFile
// Note: for Gradle, we are not returning the name unless it's a .kts file.
// This is a workaround for a project naming problem happening in Registry
// (legacy projects are named without "build.gradle" attached to them).
// See ticket BST-529 re permanent solution.
function targetFileFilteredForCompatibility(targetFile) {
    return path.basename(targetFile) === 'build.gradle.kts'
        ? targetFile
        : undefined;
}
function extractJsonFromScriptOutput(stdoutText) {
    const lines = stdoutText.split('\n');
    let jsonLine = null;
    lines.forEach((l) => {
        if (/^JSONDEPS /.test(l)) {
            if (jsonLine !== null) {
                throw new Error('More than one line with "JSONDEPS " prefix was returned; full output:\n' +
                    stdoutText);
            }
            jsonLine = l.substr(9);
        }
    });
    if (jsonLine === null) {
        throw new Error('No line prefixed with "JSONDEPS " was returned; full output:\n' +
            stdoutText);
    }
    debugLog('The command produced JSONDEPS output of ' +
        jsonLine.length +
        ' characters');
    return JSON.parse(jsonLine);
}
async function buildGraph(snykGraph, projectName, projectVersion) {
    const pkgManager = { name: 'gradle' };
    const isEmptyGraph = !snykGraph || Object.keys(snykGraph).length === 0;
    const depGraphBuilder = new dep_graph_1.DepGraphBuilder(pkgManager, {
        name: projectName,
        version: projectVersion || '0.0.0',
    });
    if (isEmptyGraph) {
        return depGraphBuilder.build();
    }
    const childrenChain = new Map();
    const ancestorsChain = new Map();
    for (const id of Object.keys(snykGraph)) {
        const { name, version, parentIds } = snykGraph[id];
        const nodeId = `${name}@${version}`;
        depGraphBuilder.addPkgNode({ name, version }, nodeId);
        if (parentIds && nodeId) {
            for (const parentId of parentIds) {
                const currentChildren = childrenChain.get(parentId) || [];
                const currentAncestors = ancestorsChain.get(parentId) || [];
                childrenChain.set(parentId, [...currentChildren, nodeId]);
                ancestorsChain.set(nodeId, [...currentAncestors, parentId]);
            }
        }
    }
    const parentlessDueOfCycles = new Set();
    // Edges
    for (const id of Object.keys(snykGraph)) {
        snykGraph[id].parentIds = Array.from(new Set(snykGraph[id].parentIds).values());
        const { name, version, parentIds } = snykGraph[id];
        const nodeId = `${name}@${version}`;
        if (parentIds && parentIds.length > 0) {
            for (let parentId of parentIds) {
                // case of missing assign version
                if (!parentId.includes('@') && snykGraph[parentId]) {
                    const { name, version } = snykGraph[parentId];
                    parentId = `${name}@${version}`;
                }
                if (!parentId || !nodeId || parentId === nodeId) {
                    continue;
                }
                const alreadyVisited = new Set();
                const hasCycles = find_cycles_1.findCycles(ancestorsChain, childrenChain, parentId, nodeId, alreadyVisited);
                if (hasCycles) {
                    parentlessDueOfCycles.add(nodeId);
                    continue;
                }
                depGraphBuilder.connectDep(parentId, nodeId);
            }
        }
        else {
            depGraphBuilder.connectDep('root-node', nodeId);
        }
    }
    //@boost temporary solution while we do not allow cycles (jvm supports cycles)
    parentlessDueOfCycles.forEach((child) => {
        depGraphBuilder.connectDep('root-node', child);
    });
    childrenChain.clear();
    ancestorsChain.clear();
    return depGraphBuilder.build();
}
exports.buildGraph = buildGraph;
async function getAllDepsOneProject(root, targetFile, options, subProject) {
    const allProjectDeps = await getAllDeps(root, targetFile, options);
    const allSubProjectNames = allProjectDeps.allSubProjectNames;
    if (subProject) {
        const { depGraph, meta } = getDepsSubProject(subProject, allProjectDeps);
        return {
            depGraph,
            allSubProjectNames,
            gradleProjectName: meta.gradleProjectName,
            versionBuildInfo: allProjectDeps.versionBuildInfo,
        };
    }
    const { projects, defaultProject } = allProjectDeps;
    const { depGraph } = projects[defaultProject];
    return {
        depGraph,
        allSubProjectNames,
        gradleProjectName: defaultProject,
        versionBuildInfo: allProjectDeps.versionBuildInfo,
    };
}
function getDepsSubProject(subProject, allProjectDeps) {
    const gradleProjectName = `${allProjectDeps.defaultProject}/${subProject}`;
    if (!allProjectDeps.projects || !allProjectDeps.projects[subProject]) {
        throw new errors_1.MissingSubProjectError(subProject, Object.keys(allProjectDeps));
    }
    const depGraph = allProjectDeps.projects[subProject].depGraph;
    return {
        depGraph,
        meta: {
            gradleProjectName,
        },
    };
}
async function getAllDepsAllProjects(root, targetFile, options) {
    const allProjectDeps = await getAllDeps(root, targetFile, options);
    return Object.keys(allProjectDeps.projects).map((proj) => {
        const defaultProject = allProjectDeps.defaultProject;
        const gradleProjectName = proj === defaultProject ? defaultProject : `${defaultProject}/${proj}`;
        return {
            targetFile: targetFileFilteredForCompatibility(allProjectDeps.projects[proj].targetFile),
            meta: {
                gradleProjectName,
                versionBuildInfo: allProjectDeps.versionBuildInfo,
                targetFile: allProjectDeps.projects[proj].targetFile,
            },
            depGraph: allProjectDeps.projects[proj].depGraph,
        };
    });
}
const reEcho = /^SNYKECHO (.*)$/;
async function printIfEcho(line) {
    const maybeMatch = reEcho.exec(line);
    if (maybeMatch) {
        debugLog(maybeMatch[1]);
    }
}
// <insert a npm left-pad joke here>
function leftPad(s, n) {
    return ' '.repeat(Math.max(n - s.length, 0)) + s;
}
async function getInjectedScriptPath() {
    let initGradleAsset = null;
    if (/index.js$/.test(__filename)) {
        // running from ./dist
        // path.join call has to be exactly in this format, needed by "pkg" to build a standalone Snyk CLI binary:
        // https://www.npmjs.com/package/pkg#detecting-assets-in-source-code
        initGradleAsset = path.join(__dirname, '../lib/init.gradle');
    }
    else if (/index.ts$/.test(__filename)) {
        // running from ./lib
        initGradleAsset = path.join(__dirname, 'init.gradle');
    }
    else {
        throw new Error('Cannot locate Snyk init.gradle script');
    }
    // We could be running from a bundled CLI generated by `pkg`.
    // The Node filesystem in that case is not real: https://github.com/zeit/pkg#snapshot-filesystem
    // Copying the injectable script into a temp file.
    try {
        const tmpInitGradle = tmp.fileSync({ postfix: '-init.gradle' });
        fs.createReadStream(initGradleAsset).pipe(fs.createWriteStream('', { fd: tmpInitGradle.fd }));
        return {
            injectedScripPath: tmpInitGradle.name,
            cleanupCallback: tmpInitGradle.removeCallback,
        };
    }
    catch (error) {
        error.message =
            error.message +
                '\n\n' +
                'Failed to create a temporary file to host Snyk init script for Gradle build analysis.';
        throw error;
    }
}
// when running a project is making use of gradle wrapper, the first time we run `gradlew -v`, the download
// process happens, cluttering the parsing of the gradle output.
// will extract the needed data using a regex
function cleanupVersionOutput(gradleVersionOutput) {
    // Everything since the first "------" line.
    // [\s\S] used instead of . as it's the easiest way to match \n too
    const matchedData = gradleVersionOutput.match(/(-{60}[\s\S]+$)/g);
    if (matchedData) {
        return matchedData[0];
    }
    debugLog('cannot parse gradle version output:' + gradleVersionOutput);
    return '';
}
function getVersionBuildInfo(gradleVersionOutput) {
    try {
        const cleanedVersionOutput = cleanupVersionOutput(gradleVersionOutput);
        if (cleanedVersionOutput !== '') {
            const gradleOutputArray = cleanedVersionOutput.split(/\r\n|\r|\n/);
            // from first 3 new lines, we get the gradle version
            const gradleVersion = gradleOutputArray[1].split(' ')[1].trim();
            const versionMetaInformation = gradleOutputArray.slice(4);
            // from line 4 until the end we get multiple meta information such as Java, Groovy, Kotlin, etc.
            const metaBuildVersion = {};
            // Select the lines in "Attribute: value format"
            versionMetaInformation
                .filter((value) => value && value.length > 0 && value.includes(': '))
                .map((value) => value.split(/(.*): (.*)/))
                .forEach((splitValue) => (metaBuildVersion[toCamelCase(splitValue[1].trim())] = splitValue[2].trim()));
            return {
                gradleVersion,
                metaBuildVersion,
            };
        }
    }
    catch (error) {
        debugLog('version build info not present, skipping ahead: ' + error);
    }
}
async function getAllDeps(root, targetFile, options) {
    const command = getCommand(root, targetFile);
    debugLog('`gradle -v` command run: ' + command);
    let gradleVersionOutput = '[COULD NOT RUN gradle -v]';
    try {
        gradleVersionOutput = await subProcess.execute(command, ['-v'], {
            cwd: root,
        });
    }
    catch (_) {
        // intentionally empty
    }
    if (gradleVersionOutput.match(/Gradle 1/)) {
        throw new Error('Gradle 1.x is not supported');
    }
    const { injectedScripPath, cleanupCallback } = await getInjectedScriptPath();
    const args = buildArgs(root, targetFile, injectedScripPath, options);
    const fullCommandText = 'gradle command: ' + command + ' ' + args.join(' ');
    debugLog('Executing ' + fullCommandText);
    try {
        const stdoutText = await subProcess.execute(command, args, { cwd: root }, printIfEcho);
        if (cleanupCallback) {
            cleanupCallback();
        }
        const extractedJson = extractJsonFromScriptOutput(stdoutText);
        const versionBuildInfo = getVersionBuildInfo(gradleVersionOutput);
        if (versionBuildInfo) {
            extractedJson.versionBuildInfo = versionBuildInfo;
        }
        // processing snykGraph from gradle task to depGraph
        for (const projectId in extractedJson.projects) {
            const { snykGraph, projectVersion } = extractedJson.projects[projectId];
            let projectName = path.basename(root);
            if (projectId !== extractedJson.defaultProject) {
                projectName = `${path.basename(root)}/${projectId}`;
            }
            extractedJson.projects[projectId].depGraph = await buildGraph(snykGraph, projectName, projectVersion);
            // this property usage ends here
            delete extractedJson.projects[projectId].snykGraph;
        }
        return extractedJson;
    }
    catch (error0) {
        const error = error0;
        const gradleErrorMarkers = /^\s*>\s.*$/;
        const gradleErrorEssence = error.message
            .split('\n')
            .filter((l) => gradleErrorMarkers.test(l))
            .join('\n');
        const orange = chalk.rgb(255, 128, 0);
        const blackOnYellow = chalk.bgYellowBright.black;
        gradleVersionOutput = orange(gradleVersionOutput);
        let mainErrorMessage = `Error running Gradle dependency analysis.

Please ensure you are calling the \`snyk\` command with correct arguments.
If the problem persists, contact support@snyk.io, providing the full error
message from above, starting with ===== DEBUG INFORMATION START =====.`;
        // Special case for Android, where merging the configurations is sometimes
        // impossible.
        // There are no automated tests for this yet (setting up Android SDK is quite problematic).
        // See test/manual/README.md
        if (cannotResolveVariantMarkers.find((m) => error.message.includes(m))) {
            // Extract attribute information via JSONATTRS marker:
            const jsonAttrs = JSON.parse(error.message
                .split('\n')
                .filter((line) => /^JSONATTRS /.test(line))[0]
                .substr(10));
            const attrNameWidth = Math.max(...Object.keys(jsonAttrs).map((name) => name.length));
            const jsonAttrsPretty = Object.keys(jsonAttrs)
                .map((name) => chalk.whiteBright(leftPad(name, attrNameWidth)) +
                ': ' +
                chalk.gray(jsonAttrs[name].join(', ')))
                .join('\n');
            mainErrorMessage = `Error running Gradle dependency analysis.

It seems like you are scanning an Android build with ambiguous dependency variants.
We cannot automatically resolve dependencies for such builds.

You have several options to make dependency resolution rules more specific:

1. Run Snyk CLI tool with an attribute filter, e.g.:
    ${chalk.whiteBright('snyk test --all-sub-projects --configuration-attributes=buildtype:release,usage:java-runtime')}

The filter will select matching attributes from those found in your configurations, use them
to select matching configuration(s) to be used to resolve dependencies. Any sub-string of the full
attribute name is enough.

Select the values for the attributes that would allow to unambiguously select the correct dependency
variant. The Gradle error message above should contain information about attributes found in
different variants.

Suggested attributes: buildtype, usage and your "flavor dimension" for Android builds.

The following attributes and their possible values were found in your configurations:
${jsonAttrsPretty}

2. Run Snyk CLI tool for specific configuration(s), e.g.:
    ${chalk.whiteBright("snyk test --gradle-sub-project=my-app --configuration-matching='^releaseRuntimeClasspath$'")}

(note that some configurations won't be present in every your subproject)

3. Converting your subproject dependency specifications from the form of
    ${chalk.whiteBright("implementation project(':mymodule')")}
to
    ${chalk.whiteBright("implementation project(path: ':mymodule', configuration: 'default')")}`;
        }
        error.message = `${chalk.red.bold('Gradle Error (short):\n' + gradleErrorEssence)}

${blackOnYellow('===== DEBUG INFORMATION START =====')}
${orange(fullCommandText)}
${orange(gradleVersionOutput)}
${orange(error.message)}
${blackOnYellow('===== DEBUG INFORMATION END =====')}

${chalk.red.bold(mainErrorMessage)}`;
        throw error;
    }
}
function toCamelCase(input) {
    input = input
        .toLowerCase()
        .replace(/(?:(^.)|([-_\s]+.))/g, (match) => {
        return match.charAt(match.length - 1).toUpperCase();
    });
    return input.charAt(0).toLowerCase() + input.substring(1);
}
function getCommand(root, targetFile) {
    const isWinLocal = /^win/.test(os.platform()); // local check, can be stubbed in tests
    const quotLocal = isWinLocal ? '"' : "'";
    const wrapperScript = isWinLocal ? 'gradlew.bat' : './gradlew';
    // try to find a sibling wrapper script first
    let pathToWrapper = path.resolve(root, path.dirname(targetFile), wrapperScript);
    if (fs.existsSync(pathToWrapper)) {
        return quotLocal + pathToWrapper + quotLocal;
    }
    // now try to find a wrapper in the root
    pathToWrapper = path.resolve(root, wrapperScript);
    if (fs.existsSync(pathToWrapper)) {
        return quotLocal + pathToWrapper + quotLocal;
    }
    return 'gradle';
}
function buildArgs(root, targetFile, initGradlePath, options) {
    const args = [];
    args.push('snykResolvedDepsJson', '-q');
    if (targetFile) {
        if (!fs.existsSync(path.resolve(root, targetFile))) {
            throw new Error('File not found: "' + targetFile + '"');
        }
        args.push('--build-file');
        let formattedTargetFile = targetFile;
        if (/\s/.test(targetFile)) {
            // checking for whitespaces
            formattedTargetFile = quot + targetFile + quot;
        }
        args.push(formattedTargetFile);
    }
    // Arguments to init script are supplied as properties: https://stackoverflow.com/a/48370451
    if (options['configuration-matching']) {
        args.push(`-Pconfiguration=${quot}${options['configuration-matching']}${quot}`);
    }
    if (options['configuration-attributes']) {
        args.push(`-PconfAttr=${quot}${options['configuration-attributes']}${quot}`);
    }
    if (!options.daemon) {
        args.push('--no-daemon');
    }
    // Parallel builds can cause race conditions and multiple JSONDEPS lines in the output
    // Gradle 4.3.0+ has `--no-parallel` flag, but we want to support older versions.
    // Not `=false` to be compatible with 3.5.x: https://github.com/gradle/gradle/issues/1827
    args.push('-Dorg.gradle.parallel=');
    // Since version 4.3.0+ Gradle uses different console output mechanism. Default mode is 'auto',
    // if Gradle is attached to a terminal. It means build output will use ANSI control characters
    // to generate the rich output, therefore JSON cannot be parsed.
    args.push('-Dorg.gradle.console=plain');
    if (!cli_interface_1.legacyPlugin.isMultiSubProject(options)) {
        args.push('-PonlySubProject=' + (options.subProject || '.'));
    }
    args.push('-I ' + initGradlePath);
    if (options.args) {
        args.push(...options.args);
    }
    // There might be a legacy --configuration option in 'args'.
    // It has been superseded by --configuration-matching option for Snyk CLI (see buildArgs),
    // but we are handling it to support the legacy setups.
    args.forEach((a, i) => {
        // Transform --configuration=foo
        args[i] = a.replace(/^--configuration[= ]([a-zA-Z_]+)/, `-Pconfiguration=${quot}^$1$$${quot}`);
        // Transform --configuration foo
        if (a === '--configuration') {
            args[i] = `-Pconfiguration=${quot}^${args[i + 1]}$${quot}`;
            args[i + 1] = '';
        }
    });
    return args;
}
exports.exportsForTests = {
    buildArgs,
    extractJsonFromScriptOutput,
    getVersionBuildInfo,
    toCamelCase,
};
//# sourceMappingURL=index.js.map