import { AnalyzedPackage, Binary, DynamicAnalysis, StaticAnalysis } from "../analyzer/types";
export declare function parseAnalysisResults(targetImage: any, analysis: StaticAnalysis | DynamicAnalysis): {
    imageId: string;
    platform: string | undefined;
    targetOS: import("../analyzer/types").OSRelease;
    type: any;
    depInfosList: AnalyzedPackage[] | Binary[];
    binaries: AnalyzedPackage[] | Binary[] | undefined;
    imageLayers: string[];
};
