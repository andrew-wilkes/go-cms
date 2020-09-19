import { DepGraphData } from '@snyk/dep-graph';
export declare const SupportFileExtensions: string[];
export interface Artifact {
    type: string;
    data: any;
    meta: {
        [key: string]: any;
    };
}
export interface ScanResult {
    artifacts: Artifact[];
    meta: {
        [key: string]: any;
    };
}
export interface Fingerprint {
    filePath: string;
    hash: string;
}
export interface Options {
    path: string;
    debug?: boolean;
}
export interface Issue {
    pkgName: string;
    pkgVersion?: string;
    issueId: string;
    fixInfo: {
        nearestFixedInVersion?: string;
    };
}
export interface IssuesData {
    [issueId: string]: {
        id: string;
        severity: string;
        title: string;
    };
}
export interface TestResult {
    issues: Issue[];
    issuesData: IssuesData;
    depGraphData: DepGraphData;
}
