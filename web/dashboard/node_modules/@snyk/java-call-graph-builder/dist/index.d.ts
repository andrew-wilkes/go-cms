import 'source-map-support/register';
import { Graph } from 'graphlib';
import { Metrics } from './metrics';
export declare function getCallGraphMvn(targetPath: string, timeout?: number): Promise<Graph>;
export declare function getClassGraphGradle(targetPath: string, timeout?: number): Promise<Graph>;
export declare function runtimeMetrics(): Metrics;
