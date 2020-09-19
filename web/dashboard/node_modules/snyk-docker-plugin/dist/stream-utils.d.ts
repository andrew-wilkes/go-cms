/// <reference types="node" />
import { Readable } from "stream";
export declare function streamToString(stream: Readable, encoding?: string): Promise<string>;
export declare function streamToBuffer(stream: Readable): Promise<Buffer>;
export declare function streamToHash(stream: Readable): Promise<string>;
export declare function streamToJson<T>(stream: Readable): Promise<T>;
