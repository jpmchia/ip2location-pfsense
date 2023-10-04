const { merge } = require('webpack-merge')
const common = require('./webpack.config.js')
const path = require('path');

module.exports = merge(common, {
    mode: 'development',
    devtool: 'eval-source-map',
    devServer: {
        static: {
            directory: path.join(__dirname, '../dist/'),
        },
        port: 9000,        
        hot: true,
    },
})