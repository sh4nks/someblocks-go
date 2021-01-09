const path = require("path");
const webpack = require("webpack");
const MiniCssExtractPlugin = require("mini-css-extract-plugin");
const { CleanWebpackPlugin } = require('clean-webpack-plugin');


// Incoperate HtmlWebpackPlugin
// Ref: https://stackoverflow.com/questions/50227120/webpack-4-split-chunks-plugin-for-multiple-entries

module.exports = {
    entry: {
        app: "./app/index.js",
    },
    output: {
        filename: "[name].js",
        path: path.resolve("../static/")
    },
    //resolve: {
    //    symlinks: false,
    //},
    plugins: [
        new CleanWebpackPlugin({
            cleanStaleWebpackAssets: true
        }),
        new MiniCssExtractPlugin({
            filename: "[name].css",
            chunkFilename: "[id].css"
        })
    ],
    optimization: {
        splitChunks: {
            maxInitialRequests: Infinity,
            minSize: 0,
            cacheGroups: {
                vendor: {
                    test(mod, chunks) {
                        // exclude anything outside of node_modules
                        if (mod.resource && !mod.resource.includes('node_modules')) {
                            return false;
                        }

                        // Exclude CSS - We already collect the CSS
                        if (mod.constructor.name === 'CssModule') {
                            return false;
                        }

                        // return all other node modules
                        return true;
                    },
                    name: 'vendors',
                    chunks: 'all',
                    enforce: true,
                },
            },
        },
    },
    module: {
        rules: [
            // Babel Loader
            {
                test: /\.js$/,
                use: {
                    loader: 'babel-loader',
                    options: {
                        presets: ['@babel/preset-env'],
                        plugins: ["@babel/plugin-syntax-dynamic-import"]
                    },
                }
            },
            // SASS / SCSS
            {
                test: /\.scss$/,
                use: [
                    MiniCssExtractPlugin.loader,
                    'css-loader',
                    'sass-loader',
                    'postcss-loader',
                ]
            },
            // CSS
            {
                test: /\.css$/,
                use: [MiniCssExtractPlugin.loader, 'css-loader']
            },
            // Images, Icons, Fonts
            {
                test: /\.(ico|jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2)(\?.*)?$/,
                use: {
                    loader: 'file-loader',
                    options: {
                        name: '[name].[ext]',
                    },
                },
            },
        ]
    }
}
