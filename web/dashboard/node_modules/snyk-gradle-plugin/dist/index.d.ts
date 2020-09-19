import { DepGraph } from '@snyk/dep-graph';
import { legacyPlugin as api } from '@snyk/cli-interface';
export interface GradleInspectOptions {
    'configuration-matching'?: string;
    'configuration-attributes'?: string;
    daemon?: boolean;
}
declare type Options = api.InspectOptions & GradleInspectOptions;
declare type VersionBuildInfo = api.VersionBuildInfo;
export declare function inspect(root: string, targetFile: string, options?: api.SingleSubprojectInspectOptions & GradleInspectOptions): Promise<api.SinglePackageResult>;
export declare function inspect(root: string, targetFile: string, options: api.MultiSubprojectInspectOptions & GradleInspectOptions): Promise<api.MultiProjectResult>;
export interface JsonDepsScriptResult {
    defaultProject: string;
    projects: ProjectsDict;
    allSubProjectNames: string[];
    versionBuildInfo: VersionBuildInfo;
}
interface SnykGraph {
    name: string;
    version: string;
    parentIds: string[];
}
interface ProjectsDict {
    [project: string]: GradleProjectInfo;
}
interface GradleProjectInfo {
    depGraph: DepGraph;
    snykGraph: {
        [name: string]: SnykGraph;
    };
    targetFile: string;
    projectVersion: string;
}
declare function extractJsonFromScriptOutput(stdoutText: string): JsonDepsScriptResult;
export declare function buildGraph(snykGraph: {
    [dependencyId: string]: SnykGraph;
}, projectName: string, projectVersion: string): Promise<DepGraph>;
declare function getVersionBuildInfo(gradleVersionOutput: string): VersionBuildInfo | undefined;
declare function toCamelCase(input: string): string;
declare function buildArgs(root: string, targetFile: string | null, initGradlePath: string, options: Options): string[];
export declare const exportsForTests: {
    buildArgs: typeof buildArgs;
    extractJsonFromScriptOutput: typeof extractJsonFromScriptOutput;
    getVersionBuildInfo: typeof getVersionBuildInfo;
    toCamelCase: typeof toCamelCase;
};
export {};
