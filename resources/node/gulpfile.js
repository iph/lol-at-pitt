var fs = require('fs');
var gulp = require('gulp');
var gutil = require('gulp-util');
var os = require('os');
var path = require('path');
var server = require('webpack-dev-server');
var webpack = require('webpack');
var webpackDevelopment = require('./webpack.development');
var webpackProduction = require('./webpack.production');

// Useful linkage to the main package directory.
var ROOT_PATH = path.resolve(__dirname);

gulp.task('build', function (callback) {

    webpack(webpackProduction, function (err, stats) {
        if (err) throw new gutil.PluginError('webpack', err);
    });

});

gulp.task('clean', function (callback) {

    var removeRecursive = function (item) {
        if (fs.existsSync(item)) {
            var stats = fs.lstatSync(item);
            if (stats.isDirectory()) {
                fs.readdirSync(item).forEach(function (child) {
                    removeRecursive(path.join(item, child));
                });
                fs.rmdirSync(item);
            } else {
                fs.unlinkSync(item);
            }
        }
    };

    // symbolic links are just deleted, not recursed
    removeRecursive(path.join(ROOT_PATH, 'dist'));
    removeRecursive(path.join(ROOT_PATH, 'build/properties'));
    removeRecursive(path.join(ROOT_PATH, 'build/webapps'));
    removeRecursive(path.join(ROOT_PATH, 'build'));

    callback();

});

gulp.task('server', function (callback) {

    // Starts up a development server to host the application assets.
    // Assets on this server will be hot loaded as changes occur.

    var HOST = os.hostname();

    var localCheck = /.*ant.*/;
    if (localCheck.test(HOST)) {
        HOST = "localhost";
    }

    var PORT = 9090;
    gutil.log('[webpack-dev-server]', ROOT_PATH);

    new server(webpack(webpackDevelopment), {

        publicPath: webpackDevelopment.output.publicPath,
        contentBase: path.join(ROOT_PATH, "entry"),
        hot: true,
        progress: true,
        headers: {
            'Access-Control-Allow-Origin': '*',
            'Access-Control-Allow-Headers': 'Origin, X-Requested-With, Content-Type, Accept'
        },
        stats: {
            colors: true
        }

    }).listen(PORT, HOST, function (err) {
        if (err) {
            throw new gutil.PluginError('webpack-dev-server', err);
        }
        gutil.log('[webpack-dev-server]', 'http://' + HOST + ':' + PORT + '/');
    });

});
