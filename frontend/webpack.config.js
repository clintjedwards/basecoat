const path = require("path");
const VueLoaderPlugin = require("vue-loader/lib/plugin");
//const CompressionPlugin = require('compression-webpack-plugin');
const webpack = require("webpack");

module.exports = {
  entry: path.resolve(__dirname, "./src/index.ts"),
  context: path.resolve(__dirname, "frontend"),
  output: {
    filename: "main.js",
    path: path.resolve(__dirname, "public/javascript"),
  },
  module: {
    rules: [
      {
        test: /\.vue$/,
        loader: "vue-loader",
      },
      {
        test: /\.tsx?$/,
        loader: "ts-loader",
        exclude: /node_modules/,
        options: {
          appendTsSuffixTo: [/\.vue$/],
        },
      },
      {
        test: /\.js$/,
        exclude: /node_modules/,
        loader: "babel-loader",
      },
      {
        test: /\.css$/,
        use: ["vue-style-loader", "css-loader"],
      },
    ],
  },
  resolve: {
    extensions: [".ts", ".js", ".vue", ".json"],
    alias: {
      vue$: "vue/dist/vue.esm.js",
    },
  },
  plugins: [
    new VueLoaderPlugin(),
    // Ignore all locale files of moment.js
    new webpack.IgnorePlugin(/^\.\/locale$/, /moment$/),
  ],
};
