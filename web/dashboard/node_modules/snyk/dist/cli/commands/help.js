"use strict";
const fs = require("fs");
const path = require("path");
const Debug = require("debug");
const debug = Debug('snyk');
module.exports = async function help(item) {
    if (!item || item === true || typeof item !== 'string') {
        item = 'help';
    }
    // cleanse the filename to only contain letters
    // aka: /\W/g but figured this was easier to read
    item = item.replace(/[^a-z-]/gi, '');
    const filename = path.resolve(__dirname, '../../../help', item + '.txt');
    try {
        return fs.readFileSync(filename, 'utf8');
    }
    catch (error) {
        debug(error);
        return `'${item}' help can't be found at location: ${filename}`;
    }
};
//# sourceMappingURL=help.js.map