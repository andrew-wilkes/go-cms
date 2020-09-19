"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TestLimitReachedError = void 0;
const custom_error_1 = require("./custom-error");
function TestLimitReachedError(errorMessage = 'Test limit reached!', errorCode = 429) {
    const error = new custom_error_1.CustomError(errorMessage);
    error.code = errorCode;
    error.strCode = 'TEST_LIMIT_REACHED';
    error.userMessage = errorMessage;
    return error;
}
exports.TestLimitReachedError = TestLimitReachedError;
//# sourceMappingURL=test-limit-reached-error.js.map