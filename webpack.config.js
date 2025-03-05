import webpack from 'webpack';
import path from 'path';

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