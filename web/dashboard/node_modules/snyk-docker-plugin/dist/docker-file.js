"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const dockerfile_ast_1 = require("dockerfile-ast");
const fs = require("fs");
const path_1 = require("path");
const instruction_parser_1 = require("./instruction-parser");
async function readDockerfileAndAnalyse(targetFilePath) {
    if (!targetFilePath) {
        return undefined;
    }
    const contents = await readFile(path_1.normalize(targetFilePath));
    return analyseDockerfile(contents);
}
exports.readDockerfileAndAnalyse = readDockerfileAndAnalyse;
async function analyseDockerfile(contents) {
    const dockerfile = dockerfile_ast_1.DockerfileParser.parse(contents);
    const baseImage = instruction_parser_1.getDockerfileBaseImageName(dockerfile);
    const dockerfilePackages = instruction_parser_1.getPackagesFromRunInstructions(dockerfile);
    const dockerfileLayers = instruction_parser_1.getDockerfileLayers(dockerfilePackages);
    return {
        baseImage,
        dockerfilePackages,
        dockerfileLayers,
    };
}
exports.analyseDockerfile = analyseDockerfile;
async function readFile(path) {
    return new Promise((resolve, reject) => {
        fs.readFile(path, "utf8", (err, data) => {
            return err ? reject(err) : resolve(data);
        });
    });
}
//# sourceMappingURL=docker-file.js.map