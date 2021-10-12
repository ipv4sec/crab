var webpack = require('webpack');
var path = require('path');


var config = {
    entry: {
        bundle_dp_channel: ['./devSrc/components/channel/index.js']
    },
    output: {
        path: path.resolve(__dirname, 'public/dist/js'),
        filename: 'bundle_dp_channel.js',

        publicPath: "",
        // chunkFilename: 'chunk_[name].js',    //TODO
    },
    // resolve: {
    //     extensions: ['.js', '.jsx']
    // },
    module: {
        rules: [
            {
                test: /\.js$/,
                loader: 'babel-loader',
                query: {
                    presets: ['es2015', 'stage-0', 'react' ],
                    plugins: ['babel-plugin-transform-runtime']
                },
                exclude: /node_modules/,
            },
           
        ]
    }
};

module.exports = config;