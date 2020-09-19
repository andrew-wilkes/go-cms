/// <reference types="node" />
import { Readable } from "stream";
export declare type ExtractCallback = (dataStream: Readable) => Promise<string | Buffer>;
export declare type FileNameAndContent = Record<string, string | Buffer>;
export interface ExtractionResult {
    imageId: string;
    manifestLayers: string[];
    extractedLayers: ExtractedLayers;
    rootFsLayers?: string[];
    platform?: string;
}
export interface ExtractedLayers {
    [layerName: string]: FileNameAndContent;
}
export interface ExtractedLayersAndManifest {
    layers: ExtractedLayers[];
    manifest: DockerArchiveManifest;
    imageConfig: DockerArchiveImageConfig;
}
export interface DockerArchiveManifest {
    Config: string;
    RepoTags: string[];
    Layers: string[];
}
export interface DockerArchiveImageConfig {
    architecture: string;
    os: string;
    rootfs: {
        diff_ids: string[];
    };
}
export interface OciArchiveLayer {
    digest: string;
}
export interface OciArchiveManifest {
    schemaVersion: string;
    config: {
        digest: string;
    };
    layers: OciArchiveLayer[];
}
export interface OciManifestInfo {
    digest: string;
    platform?: {
        architecture: string;
        os: string;
    };
}
export interface OciImageIndex {
    manifests: OciManifestInfo[];
}
export interface ExtractAction {
    actionName: string;
    filePathMatches: (filePath: string) => boolean;
    callback?: ExtractCallback;
}
