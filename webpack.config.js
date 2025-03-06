import webpack from 'webpack';
import path from 'path';
import LicensePlugin from 'webpack-license-plugin'

export default {
    entry: "./index.js",
    output: {
        filename: `built.js`,
    },
    target: ['node', 'es5'],
    mode: 'production',
    devtool: false,
    plugins: [
        new webpack.optimize.LimitChunkCountPlugin({
            maxChunks: 1
        }),
        new LicensePlugin({
            additionalFiles: {
                '../NOTICE': (pkgs) => {
                    let res = "NOTICE\n"
                    res += "see dist/oss-licenses.json and dist/built.js.LICENSE.txt for more information"
                    res += "below are dependencies including their author and the corresponding license that are (partially) included in the dist/built.js as compiled by Webpack\n"
                    res += "\n"
                    for (let pkg of pkgs) {
                        res += pkg['name']
                        if (pkg['author']) {
                            res += ` (by ${pkg['author']})`
                        }
                        res += `: ${pkg['license']}\n`
                        res += `repository: ${pkg['repository']}\n`
                        res += "\n"
                    }
                    return res
                }
            },
            includeNoticeText: true,
            unacceptableLicenseTest: (licenseIdentifier) => {
                // only allow the following licenses
                return !['MIT', 'Apache-2.0', 'ISC', 'BSD-3-Clause', 'BSD-2-Clause', '0BSD'].includes(licenseIdentifier)
            }
        })
    ],
    optimization: {
        minimize: true,
    },
    resolve: {
        alias: {
            // --watch is not supported hence fsevents.node is not required
            [path.resolve("./node_modules/fsevents/fsevents.node")]:
                path.resolve("node/fsevents/index.js")
        }
    },
    module: {
        rules: [
            {
                test: /.node$/,
                loader: 'node-loader',
            },
            {
                // Only run `.js` files through Babel
                test: /\.m?js$/,
                // exclude: /(node_modules)/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env'],
                        compact: true
                    }
                }
            },
            {
                test: /compile\/codegen\/code\.js$/,
                loader: 'string-replace-loader',
                options: {
                    search: /\s{4}return JSON.stringify\(x\)\n\s{8}.replace\(\/\\u2028\/g, "\\\\u2028"\)\n\s{8}.replace\(\/\\u2029\/g, "\\\\u2029"\);/,
                    replace() {
                        return "    return JSON.stringify(x);"
                    },
                    flags: 'g',
                    strict: true,
                }
            }
        ]
    }
};