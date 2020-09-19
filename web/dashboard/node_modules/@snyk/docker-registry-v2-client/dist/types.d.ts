export interface LayerConfig {
    mediaType: string;
    size: number;
    digest: string;
}
export interface ImageManifest {
    schemaVersion: number;
    mediaType: string;
    config: LayerConfig;
    layers: LayerConfig[];
}
