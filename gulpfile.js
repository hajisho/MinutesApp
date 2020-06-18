const { src, dest, series } = require("gulp");
const del = require('del');
const webpackStream = require("webpack-stream");
const webpack = require("webpack");

const webpackConfig = require("./webpack.config");

function clean() {
    return del(["dist/**"]);
}

function use_webpack(){
    return webpackStream(webpackConfig, webpack)
      .pipe(dest("dist"));
}

function copy() {
    return src(["src/public/**/*.html","src/public/**/*.css"])
        .pipe(dest("dist/public"));
}

exports.default = series(clean, use_webpack, copy);
