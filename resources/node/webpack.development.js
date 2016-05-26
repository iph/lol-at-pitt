var path = require('path');
var os = require('os');
var webpack = require('webpack');

var PORT = 9090;
var ROOT_PATH = path.resolve(__dirname);


var HOST = os.hostname();
var localCheck = /.*ant.*/;
if (localCheck.test(HOST)) {
    HOST = "localhost";
}

var SERVER = 'http://' + HOST + ':' + PORT;


// Grab the node path environment variable, injected by Brazil
// Split on the seperate and pass as array into the path resolver.
// Expected format: "NODE_PATH: :PATH1:PATH2"
var modulePaths = path.join(ROOT_PATH, "node_modules");
if (process.env.NODE_PATH) {
    modulePaths = process.env.NODE_PATH.slice(1).split(':');
}

module.exports = {

    // Generate a source map.
    devtool: 'source-map',

    // Load the development server, hot loader, and hub entry point.
    entry: {
        draft: [
            'webpack-dev-server/client?' + SERVER,
            'webpack/hot/only-dev-server',
            path.join(ROOT_PATH, 'app/draft.js')
        ]
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
        filename: '[name].js',
        publicPath: '/public/'
    },

    module: {
        preLoaders: [],
        loaders: [
            {
                test: /(\.js|\.jsx)$/,
                loader: 'react-hot',
                exclude: [/commonjs/, /node_modules/]
            },
            {
                test: /(\.js|\.jsx)$/,
                loader: 'babel',
                query: {
                    presets: [
                        "es2015",
                        "react"
                    ]
                },
                exclude: [/commonjs/, /node_modules/]
            },
            {
                test: /\.css$/,
                loaders: ['style', 'css'],
                exclude: /commonjs/
            }
        ]
    },

    plugins: [
        new webpack.HotModuleReplacementPlugin()
    ]

};
