"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const _isEmpty = require("lodash.isempty");
const _set = require("lodash.set");
const pMap = require("p-map");
const _1 = require("./");
const event_loop_spinner_1 = require("event-loop-spinner");
const errors_1 = require("../errors");
const config_1 = require("../config");
const EVENT_PROCESSING_CONCURRENCY = 5;
class YarnLockParseBase {
    constructor(type) {
        this.type = type;
        this.eventLoopSpinRate = 20;
        // Number of dependencies including root one.
        this.treeSize = 1;
    }
    async getDependencyTree(manifestFile, lockfile, includeDev = false, strict = true) {
        var _a;
        if (lockfile.type !== this.type) {
            throw new errors_1.InvalidUserInputError('Unsupported lockfile provided. Please provide `yarn.lock`.');
        }
        const yarnLock = lockfile;
        const depTree = {
            dependencies: {},
            hasDevDependencies: !_isEmpty(manifestFile.devDependencies),
            name: manifestFile.name,
            size: 1,
            version: manifestFile.version || '',
        };
        const nodeVersion = (_a = manifestFile === null || manifestFile === void 0 ? void 0 : manifestFile.engines) === null || _a === void 0 ? void 0 : _a.node;
        if (nodeVersion) {
            _set(depTree, 'meta.nodeVersion', nodeVersion);
        }
        const packageManagerVersion = lockfile.type === _1.LockfileType.yarn ? '1' : '2';
        _set(depTree, 'meta.packageManagerVersion', packageManagerVersion);
        const topLevelDeps = _1.getTopLevelDeps(manifestFile, includeDev);
        // asked to process empty deps
        if (_isEmpty(manifestFile.dependencies) && !includeDev) {
            return depTree;
        }
        await pMap(topLevelDeps, (dep) => this.resolveDep(dep, depTree, yarnLock, strict), { concurrency: EVENT_PROCESSING_CONCURRENCY });
        depTree.size = this.treeSize;
        return depTree;
    }
    async buildSubTree(lockFile, tree, strict) {
        const queue = [{ path: [], tree }];
        while (queue.length > 0) {
            const queueItem = queue.pop();
            const depKey = `${queueItem.tree.name}@${queueItem.tree.version}`;
            const dependency = lockFile.object[depKey];
            if (!dependency) {
                if (strict) {
                    throw new errors_1.OutOfSyncError(queueItem.tree.name, this.type);
                }
                if (!queueItem.tree.labels) {
                    queueItem.tree.labels = {};
                }
                queueItem.tree.labels.missingLockFileEntry = 'true';
                continue;
            }
            // Overwrite version pattern with exact version.
            queueItem.tree.version = dependency.version;
            if (queueItem.path.indexOf(depKey) >= 0) {
                if (!queueItem.tree.labels) {
                    queueItem.tree.labels = {};
                }
                queueItem.tree.labels.pruned = 'cyclic';
                continue;
            }
            const subDependencies = Object.entries(Object.assign(Object.assign({}, dependency.dependencies), dependency.optionalDependencies));
            for (const [subName, subVersion] of subDependencies) {
                // tree size limit should be 6 millions.
                if (this.treeSize > config_1.config.YARN_TREE_SIZE_LIMIT) {
                    throw new errors_1.TreeSizeLimitError();
                }
                const subDependency = {
                    labels: {
                        scope: tree.labels.scope,
                    },
                    name: subName,
                    version: subVersion,
                };
                if (!queueItem.tree.dependencies) {
                    queueItem.tree.dependencies = {};
                }
                queueItem.tree.dependencies[subName] = subDependency;
                queue.push({
                    path: [...queueItem.path, depKey],
                    tree: subDependency,
                });
                this.treeSize++;
                if (event_loop_spinner_1.eventLoopSpinner.isStarving()) {
                    await event_loop_spinner_1.eventLoopSpinner.spin();
                }
            }
        }
        return tree;
    }
    async resolveDep(dep, depTree, yarnLock, strict) {
        if (/^file:/.test(dep.version)) {
            depTree.dependencies[dep.name] = _1.createDepTreeDepFromDep(dep);
        }
        else {
            depTree.dependencies[dep.name] = await this.buildSubTree(yarnLock, _1.createDepTreeDepFromDep(dep), strict);
        }
        this.treeSize++;
        if (event_loop_spinner_1.eventLoopSpinner.isStarving()) {
            await event_loop_spinner_1.eventLoopSpinner.spin();
        }
    }
}
exports.YarnLockParseBase = YarnLockParseBase;
//# sourceMappingURL=yarn-lock-parse-base.js.map