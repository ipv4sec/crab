const path = require('path')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const webpack = require('webpack')

let isDev = process.env.NODE_ENV !== 'production' ? true : false

let config = {
    mode: process.env.NODE_ENV || 'none',

    entry: {
        app: path.resolve(__dirname, '../src/index.js')
    },
    output: {
        path: path.resolve(__dirname, '../../public'),
        filename: 'js/[name].[contenthash].bundle.js',
        chunkFilename: 'js/[id].chunk.js',
        publicPath: '/',
        clean: {
            keep(assets){
                console.log('clean assets',assets)
                if(/(static|images|font)\/*/.test(assets)) {
                    return true
                }
            }
        }
    },
    module: {
        rules: [
            {
                test: /\.(js|jsx)$/,
                include: path.resolve(__dirname, '../src'),
                loader: 'babel-loader'
            },
            {
                test: /\.css$/,
                use: [isDev ? 'style-loader' : MiniCssExtractPlugin.loader, 'css-loader']
            },
            {
                test: /\.(scss|sass)$/,
                use: [isDev ? 'style-loader' : MiniCssExtractPlugin.loader, 'css-loader', 'sass-loader']
            }
        ]
    },
    resolve: {
        extensions: [".jsx", '.js', '...'],
        modules: ['../../node_modules']
    },
    plugins: [
        new HtmlWebpackPlugin({
            filename: 'index.html',
            template: path.resolve(__dirname, '../src/index.html')
        }),
        new MiniCssExtractPlugin({
            filename: 'css/[name].css'
        })
    ]
}


if(isDev) {
    config.devtool = 'eval-cheap-source-map'
    config.plugins.push(new webpack.HotModuleReplacementPlugin())
    config.devServer = {
        contentBase: path.resolve(__dirname, '../../public'),
        host: '0.0.0.0',
        port: 3001,
        historyApiFallback: true,
        hot: true,
        // hotOnly: true,
        proxy: {
            '/api': 'http://localhost:3000'
        }
    }
}else {
    config.optimization = {
        minimize: true, // 默认就是true，使用terserplugin来压缩文件
        runtimeChunk: {
            name: 'runtime'
        },
        splitChunks: {
            chunks: 'async',
            minSize: 20000,
            minRemainingSize: 0,
            minChunks: 1,
            maxAsyncRequests: 30,
            maxInitialRequests: 30,
            enforceSizeThreshold: 50000,
            cacheGroups: {
              defaultVendors: {
                test: /[\\/]node_modules[\\/]/,
                priority: -10,
                reuseExistingChunk: true,
              },
              default: {
                minChunks: 2,
                priority: -20,
                reuseExistingChunk: true,
              },
            },
        }
    }

}




module.exports = config