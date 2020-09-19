import { NeedleResponse } from "needle";
export declare function registryV2Call(registryBase: string, endpoint: string, accept: string, username?: string, password?: string, reqOptions?: any): Promise<NeedleResponse>;
export declare function paginatedV2Call(registryBase: string, accept: string, username: string, password: string, endpoint: string, key: string, pageSize?: number, maxPages?: number, reqOptions?: any): Promise<string[]>;
export declare function getChallengeURL(challenge: string, registryBase: string): string;
