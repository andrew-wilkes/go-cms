"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.countPathsToGraphRoot = void 0;
function countPathsToGraphRoot(graph) {
    return graph
        .getPkgs()
        .reduce((acc, pkg) => acc + graph.countPathsToRoot(pkg), 0);
}
exports.countPathsToGraphRoot = countPathsToGraphRoot;
//# sourceMappingURL=utils.js.map