var path = require('path');
var os = require('os');
var webpack = require('webpack');

var ROOT_PATH = path.resolve(__dirname);

// Grab the node path environment variable, injected by Brazil
// Split on the seperate and pass as array into the path resolver.
// Expected format: "NODE_PATH: :PATH1:PATH2"
var modulePaths = path.join(ROOT_PATH, "node_modules/");
if (process.env.NODE_PATH) {
    modulePaths = process.env.NODE_PATH.slice(1).split(':');
}


module.exports = {

    // Generate a source map.
    devtool: null,

    // Load the development server, hot loader, and hub entry point.
    entry: {
        draft: path.join(ROOT_PATH, 'app/draft.js')
    },

    // We need to map module root to Brazil instead of the default NPM.
    resolve: {
        root: modulePaths,
        extensions: ['', '.js', '.jsx']
    },

    // We need to map loader root to Brazil instead of the default NPM.
    resolveLoader: {
        root: modulePaths
    },

    // Save processed files to the 'dist' folder.
    output: {
        path: path.join(ROOT_PATH, '../public/js'),
        filename: '[name].js'
    },

    module: {
        preLoaders: [],
        loaders: [
            {
                test: /\.js$/,
                loaders: ['babel?presets[]=es2015&presets[]=react'],
                exclude: [/commonjs/, /node_modules/]
            },
            {test: /\.css$/, loaders: ['style', 'css'], exclude: /commonjs/},
        ]
    },

    plugins: []

};
