"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const snyk_docker_pull_1 = require("@snyk/snyk-docker-pull");
const Debug = require("debug");
const Modem = require("docker-modem");
const event_loop_spinner_1 = require("event-loop-spinner");
const fs_1 = require("fs");
const minimatch = require("minimatch");
const os_1 = require("os");
const fspath = require("path");
const lsu = require("./ls-utils");
const subProcess = require("./sub-process");
const SystemDirectories = ["dev", "proc", "sys"];
const debug = Debug("snyk");
class Docker {
    constructor(targetImage, options) {
        var _a;
        this.targetImage = targetImage;
        this.optionsList = Docker.createOptionsList(options);
        this.socketPath =
            ((_a = options) === null || _a === void 0 ? void 0 : _a.socketPath) ||
                (os_1.platform() === "win32"
                    ? "\\\\.\\pipe\\docker_engine"
                    : "/var/run/docker.sock");
    }
    static async binaryExists() {
        try {
            await subProcess.execute("docker", ["version"]);
            return true;
        }
        catch (e) {
            return false;
        }
    }
    static run(args, options) {
        return subProcess.execute("docker", [
            ...Docker.createOptionsList(options),
            ...args,
        ]);
    }
    static createOptionsList(options) {
        const opts = [];
        if (!options) {
            return opts;
        }
        if (options.host) {
            opts.push(`--host=${options.host}`);
        }
        if (options.tlscert) {
            opts.push(`--tlscert=${options.tlscert}`);
        }
        if (options.tlscacert) {
            opts.push(`--tlscacert=${options.tlscacert}`);
        }
        if (options.tlskey) {
            opts.push(`--tlskey=${options.tlskey}`);
        }
        if (options.tlsverify) {
            opts.push(`--tlsverify=${options.tlsverify}`);
        }
        return opts;
    }
    /**
     * Runs the command, catching any expected errors and returning them as normal
     * stderr/stdout result.
     */
    async runSafe(cmd, args = [], 
    // no error is thrown if any of listed errors is found in stderr
    ignoreErrors = ["No such file", "not found"]) {
        try {
            return await this.run(cmd, args);
        }
        catch (error) {
            const stderr = error.stderr;
            if (typeof stderr === "string") {
                if (ignoreErrors.some((errMsg) => stderr.indexOf(errMsg) >= 0)) {
                    return { stdout: error.stdout, stderr };
                }
            }
            throw error;
        }
    }
    run(cmd, args = []) {
        return subProcess.execute("docker", [
            ...this.optionsList,
            "run",
            "--rm",
            "--entrypoint",
            '""',
            "--network",
            "none",
            this.targetImage,
            cmd,
            ...args,
        ]);
    }
    async pull(registry, repo, tag, imageSavePath, username, password) {
        const dockerPull = new snyk_docker_pull_1.DockerPull();
        const opt = {
            username,
            password,
            loadImage: false,
            imageSavePath,
        };
        return await dockerPull.pull(registry, repo, tag, opt);
    }
    async pullCli(targetImage, options) {
        var _a;
        return subProcess.execute("docker", [
            "pull",
            ((_a = options) === null || _a === void 0 ? void 0 : _a.platform) ? `--platform=${options.platform}` : "",
            targetImage,
        ]);
    }
    async save(targetImage, destination) {
        const request = {
            path: `/images/${targetImage}/get?`,
            method: "GET",
            isStream: true,
            statusCodes: {
                200: true,
                400: "bad request",
                404: "not found",
                500: "server error",
            },
        };
        debug(`Docker.save: targetImage: ${targetImage}, destination: ${destination}`);
        const modem = new Modem({ socketPath: this.socketPath });
        return new Promise((resolve, reject) => {
            modem.dial(request, (err, stream) => {
                if (err) {
                    return reject(err);
                }
                const writeStream = fs_1.createWriteStream(destination);
                writeStream.on("error", (err) => {
                    reject(err);
                });
                writeStream.on("finish", () => {
                    resolve();
                });
                stream.on("error", (err) => {
                    reject(err);
                });
                stream.on("end", () => {
                    writeStream.end();
                });
                stream.pipe(writeStream);
            });
        });
    }
    async inspectImage(targetImage) {
        return subProcess.execute("docker", [
            ...this.optionsList,
            "inspect",
            targetImage,
        ]);
    }
    async catSafe(filename) {
        return this.runSafe("cat", [filename]);
    }
    async lsSafe(path, recursive) {
        let params = "-1ap";
        if (recursive) {
            params += "R";
        }
        const ignoreErrors = [
            "No such file",
            "file not found",
            "Permission denied",
        ];
        return this.runSafe("ls", [params, path], ignoreErrors);
    }
    /**
     * Find files on a docker image according to a given list of glob expressions.
     */
    async findGlobs(globs, exclusionGlobs = [], path = "/", recursive = true, excludeRootDirectories = SystemDirectories) {
        let root;
        const res = [];
        if (recursive && path === "/") {
            // When scanning from the root of a docker image we need to
            // exclude system files e.g. /proc, /sys, etc. to make the
            // operation less expensive.
            const outputRoot = await this.lsSafe("/", false);
            root = lsu.parseLsOutput(outputRoot.stdout);
            for (const subdir of root.subDirs) {
                if (excludeRootDirectories.includes(subdir.name)) {
                    continue;
                }
                const subdirOutput = await this.lsSafe("/" + subdir.name, true);
                const subdirRecursive = lsu.parseLsOutput(subdirOutput.stdout);
                await lsu.iterateFiles(subdirRecursive, (f) => {
                    f.path = "/" + subdir.name + f.path;
                });
                subdir.subDirs = subdirRecursive.subDirs;
                subdir.files = subdirRecursive.files;
            }
        }
        else {
            const output = await this.lsSafe(path, recursive);
            if (event_loop_spinner_1.eventLoopSpinner.isStarving()) {
                await event_loop_spinner_1.eventLoopSpinner.spin();
            }
            root = lsu.parseLsOutput(output.stdout);
        }
        await lsu.iterateFiles(root, (f) => {
            const filepath = fspath.join(f.path, f.name);
            let exclude = false;
            for (const g of exclusionGlobs) {
                if (!exclude && minimatch(filepath, g)) {
                    exclude = true;
                }
            }
            if (!exclude) {
                for (const g of globs) {
                    if (minimatch(filepath, g)) {
                        res.push(filepath);
                    }
                }
            }
        });
        return res;
    }
}
exports.Docker = Docker;
//# sourceMappingURL=docker.js.map