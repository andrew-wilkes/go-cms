import 'source-map-support/register';
import { Graph } from 'graphlib';
export declare function getCallGraphGenCommandArgs(classPath: string, jarPath: string, targets: string[]): string[];
export declare function getTargets(targetPath: string): Promise<string[]>;
export declare function getClassPerJarMapping(classPath: string): Promise<{
    [index: string]: string;
}>;
export declare function getCallGraph(classPath: string, targetPath: string, timeout?: number): Promise<Graph>;
