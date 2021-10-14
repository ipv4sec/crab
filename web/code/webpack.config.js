var webpack = require('webpack');
var path = require('path');
var ExtractTextPlugin = require("extract-text-webpack-plugin");
var { CleanWebpackPlugin } = require('clean-webpack-plugin');

var config = {
    entry: {
        main: ['./src/index.js'],
        // vendor: ['react', 'react-dom', 'redux', 'material-ui','react-router', 'react-redux', 'echarts'] //Object.keys(packageJson.dependencies), 'echarts'
    },
    output: {
        path: path.resolve(__dirname, 'public/dist/js'),
        filename: 'bundle_[name].js',
        publicPath: "/dist/js/",
        chunkFilename: 'chunk_[name]_[chunkhash:8].js',
    },
    resolve: {
        extensions: ['.js', '.jsx']
    },
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
            {test: /\.css$/, loader: ExtractTextPlugin.extract("style-loader", "css-loader")},
            {
                test: /\.(gif|jpg|svg|png|woff|eot|ttf)\??.*$/,
                loader: 'url-loader?limit=8192&name=[path][name].[ext]'
            }
           
        ]
    },
    plugins: [
        new ExtractTextPlugin("component.css"),
        new webpack.optimize.CommonsChunkPlugin({
            name: 'vendor',
            filename: '[name].js'
        }),

        new webpack.DefinePlugin({
            'process.env': {
                'NODE_ENV': JSON.stringify('production')
                // 'NODE_ENV': JSON.stringify("development")
            }
        }),
        new CleanWebpackPlugin(),
        // new webpack.optimize.UglifyJsPlugin({
        //     compress:{
        //         warnings: false
        //     }
        // })
    ]
};

module.exports = config;