var gulp = require('gulp');
var gutil = require('gulp-util');
var browserify = require('gulp-browserify');
var concat = require('gulp-concat');
var sass = require('gulp-sass');
var coffeeify = require('coffeeify');

var coffeePaths = ['static/scripts/**/*.coffee'];
var mainCoffeePath = ['static/scripts/main.coffee'];
gulp.task('compile-coffee', function() {
  return gulp
    .src(mainCoffeePath, {read: false})
    .pipe(browserify({
      transform: ['coffeeify'],
      extensions: ['.coffee'],
      debug: !!"yes, I want source maps"
    }))
      .on('error', gutil.log)
    .pipe(concat('scripts.gen.js'))
    .pipe(gulp.dest('static'));
});

var sassPaths = ['static/styles/**/*.scss'];
gulp.task('compile-sass', function() {
  return gulp
    .src(sassPaths)
    .pipe(sass())
      .on('error', gutil.log)
    .pipe(concat('styles.gen.css'))
    .pipe(gulp.dest('static'));
});

gulp.task('watch-all', function() {
  gulp.watch(coffeePaths, ['compile-coffee']);
  gulp.watch(sassPaths, ['compile-sass']);
});

gulp.task('default', [
  'compile-coffee',
  'compile-sass',
  'watch-all'
]);
