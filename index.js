import {lint} from '@stoplight/spectral-cli/dist/services/linter/linter.js'
import {formatOutput} from "@stoplight/spectral-cli/dist/services/output.js";

exports.formatOutput = formatOutput
exports.lint = lint(lintDocuments, {
    encoding: "utf8",
    format: ["json"],
    output: {
        "json": "<stdout>",
    },
    ruleset: lintRuleset,
    stdinFilepath: null,
    ignoreUnknownFormat: true,
    failOnUnmatchedGlobs: false,
    verbose: true,
    quiet: false,
})