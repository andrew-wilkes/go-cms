import { Graph } from 'graphlib';
export declare function buildCallGraph(input: string, classPerJarMapping: {
    [index: string]: string;
}): Graph;
