angular = require 'angular'

exports = angular.module('jbitor.twodistances', [])

exports.factory 'twodistances', -> require './twodistances'
