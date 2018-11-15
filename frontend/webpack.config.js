const path = require('path');
const VueLoaderPlugin = require('vue-loader/lib/plugin');

module.exports = {
    entry: path.resolve(__dirname, 'src'),
    context: path.resolve(__dirname, 'frontend'),
    output: {
        filename: 'main.js',
        path: path.resolve(__dirname, 'public/javascript')
    },
    resolve: {
        alias: {
            'vue$': 'vue/dist/vue.esm.js'
        },
    },
    module: {
        rules: [{
            test: /\.vue$/,
            loader: 'vue-loader',
        },
        {
            test: /\.js$/,
            exclude: /node_modules/,
            loader: 'babel-loader'
        },
        {
            test: /\.css$/,
            use: [
                'vue-style-loader',
                'css-loader'
            ]
        }
        ]
    },
    plugins: [
        new VueLoaderPlugin()
    ]
};
