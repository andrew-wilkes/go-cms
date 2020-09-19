export interface DirResult {
  name: string;
  removeCallback: () => void;
}

export interface DockerPullResult {
  imageDigest: string;
  stagingDir: DirResult | null;
  /** @deprecated caching is no longer used */
  cachedLayersDigests: string[];
  missingLayersDigests: string[];
  pullDuration: number;
}

export interface DockerPullOptions {
  username?: string;
  password?: string;
  // weak typing on the client
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  reqOptions?: any;
  /**
   * loadImage will default to true if no value is sent
   */
  loadImage?: boolean;
  imageSavePath?: string;
}

export interface SaveRequests {
  [name: string]: SaveRequest;
}

interface SaveRequest {
  username?: string;
  registryBase?: string;
  repo?: string;
  tag?: string;
}
