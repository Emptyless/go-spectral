var spectral = require('./dist/built.js');
new Promise(function(res, rej) {
    spectral.lint
        .then(function(output) {
            res(spectral.formatOutput(output.results, "json", { failSeverity: -1 }, output.resolvedRuleset));
        })
        .catch(function(e) {
            rej(e);
        })
})