"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const yarnLockfileParser = require("@yarnpkg/lockfile");
const _1 = require("./");
const errors_1 = require("../errors");
const yarn_lock_parse_base_1 = require("./yarn-lock-parse-base");
class YarnLockParser extends yarn_lock_parse_base_1.YarnLockParseBase {
    constructor() {
        super(_1.LockfileType.yarn);
    }
    parseLockFile(lockFileContents) {
        try {
            const yarnLock = yarnLockfileParser.parse(lockFileContents);
            yarnLock.dependencies = yarnLock.object;
            yarnLock.type = _1.LockfileType.yarn;
            return yarnLock;
        }
        catch (e) {
            throw new errors_1.InvalidUserInputError(`yarn.lock parsing failed with an error: ${e.message}`);
        }
    }
}
exports.YarnLockParser = YarnLockParser;
//# sourceMappingURL=yarn-lock-parse.js.map