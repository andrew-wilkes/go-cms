import { LockfileParser, PkgTree, ManifestFile, Lockfile, LockfileType } from './';
export declare type YarnLockFileTypes = LockfileType.yarn | LockfileType.yarn2;
export interface YarnLockDeps {
    [depName: string]: YarnLockDep;
}
export interface YarnLockBase<T extends YarnLockFileTypes> {
    type: string;
    object: YarnLockDeps;
    dependencies?: YarnLockDeps;
    lockfileType: T;
}
export interface YarnLockDep {
    version: string;
    dependencies?: {
        [depName: string]: string;
    };
    optionalDependencies?: {
        [depName: string]: string;
    };
}
export declare abstract class YarnLockParseBase<T extends YarnLockFileTypes> implements LockfileParser {
    private type;
    private treeSize;
    private eventLoopSpinRate;
    constructor(type: T);
    abstract parseLockFile(lockFileContents: string): Lockfile;
    getDependencyTree(manifestFile: ManifestFile, lockfile: Lockfile, includeDev?: boolean, strict?: boolean): Promise<PkgTree>;
    private buildSubTree;
    private resolveDep;
}
