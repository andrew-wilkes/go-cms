import { DockerArchiveImageConfig, DockerArchiveManifest } from "../types";
export { extractArchive } from "./layer";
export declare function getManifestLayers(manifest: DockerArchiveManifest): string[];
export declare function getImageIdFromManifest(manifest: DockerArchiveManifest): string;
export declare function getRootFsLayersFromConfig(imageConfig: DockerArchiveImageConfig): string[];
export declare function getPlatformFromConfig(imageConfig: DockerArchiveImageConfig): string | undefined;
