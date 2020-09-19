import { LockfileType } from './';
import { YarnLockBase } from './yarn-lock-parse-base';
import { YarnLockParseBase } from './yarn-lock-parse-base';
export declare type YarnLock = YarnLockBase<LockfileType.yarn>;
export declare class YarnLockParser extends YarnLockParseBase<LockfileType.yarn> {
    constructor();
    parseLockFile(lockFileContents: string): YarnLock;
}
