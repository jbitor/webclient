jQuery = require 'jquery'

require './jbitor.twodistances'

exports = angular.module('jbitor.app', [
  'jbitor.twodistances'
])

# We define a new HREF whitelist which includes magnet URLs.
exports.config (
  $compileProvider
) ->
  $compileProvider.aHrefSanitizationWhitelist(
    /^\s*(https?|ftp|mailto|magnet):/);

exports.controller 'jbitorAppController', (
  $scope,
  $interval,
  DHTFindPeersModel,
  twodistances
) ->
  $scope.connectionError = null
  $scope.clientState =
    dht:
      connectionInfo: {}
      peerSearches: []

  $scope.$watch 'clientState.dht', (newValue) ->
    return @ unless newValue?

    # HACK: using jQuery

    console.log(newValue)

    jQuery('canvas').remove()

    for search in newValue.peerSearches
      nodes = for id, info of search.queriedNodes
        new twodistances.Node(0.49 + Math.random()*0.51, 0.49 + Math.random()*0.51, null)

      graph = new twodistances.Graph(1.0, nodes)
      graph.canvas.style.display = 'block'
      jQuery('body').prepend(graph.canvas)

    @


  # Constantly reload the clientState from the server.
  intervalPromise = $interval ->
    jQuery.getJSON('/api/clientState.json').then (state) ->
      console.log "Got state", state
      $scope.$apply ->
        $scope.connectionError = null;
        $scope.clientState = state;
        @
    , ->
      $scope.$apply ->
        $scope.connectionError = true;
        console.error "Failed to get update client state", err
        @
  , 750

  $scope.$on '$destroy', ->
    $interval.cancel intervalPromise
    @

  $scope.dht =
    connectionInfo: ->
      $scope.clientState.dht.connectionInfo

    peerSearches: ->
      $scope.clientState.dht.peerSearches

    health: ->
      if @connectionInfo().GoodNodes
        Math.min(
          1.0,
          Math.max(
            0.1,
            this.connectionInfo().GoodNodes / 32
          )
        )
      else if this.connectionInfo().UnknownNodes
        0.05

    findPeers: new DHTFindPeersModel()

  return @

exports.factory 'DHTFindPeersModel', ->
    class DHTFindPeersModel
      constructor: ->
        @infohash = 'e3811b9539cacff680e418124272177c47477157';
        @requests = []

      onSubmit: (event) ->
        jQuery.post '/api/peerSearch', infohash: this.infohash

        @infohash = ''
        event.preventDefault()
        @

    DHTFindPeersModel
