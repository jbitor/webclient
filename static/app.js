'use strict';

angular.module('jbitor.app', [])

// We define a new HREF whitelist which includes magnet URLs.
.config(function(
    $compileProvider
) {   
    $compileProvider.aHrefSanitizationWhitelist(
        /^\s*(https?|ftp|mailto|magnet):/);
})

.controller('jbitorAppController', function(
    $scope,
    $interval,
    DHTFindPeersModel
) {
    $scope.connectionError = null;
    $scope.clientState = {
        dht: {
            connectionInfo: {},
            peerSearches: []
        }
    };

    // Constantly reload the clientState from the server.
    var intervalPromise = $interval(function() {
        $.getJSON('/api/clientState.json').then(function(state) {
            $scope.$apply(function() {
                $scope.connectionError = null;
                $scope.clientState = state;
            });
        }, function(err) {
            $scope.$apply(function() {
                $scope.connectionError = true;
                console.error("Failed to get update client state", err);
            });
        });
    }, 750);

    $scope.$on('$destroy', function() {
        $interval.cancel(intervalPromise);
    });

    $scope.dht = {
        connectionInfo: function() {
            return $scope.clientState.dht.connectionInfo;
        },
        peerSearches: function() {
            return $scope.clientState.dht.peerSearches;
        },
        health: function() {
            if (this.connectionInfo().GoodNodes) {
                return Math.min(
                    1.0,
                    Math.max(
                        0.1,
                        this.connectionInfo().GoodNodes / 32
                    )
                );
            } else if (this.connectionInfo().UnknownNodes) {
                return 0.05;
            }
        },
        findPeers: new DHTFindPeersModel()
    };
})

.factory('DHTFindPeersModel', function() {
    function DHTFindPeersModel() {
        this.infohash = 'e3811b9539cacff680e418124272177c47477157';

        this.requests = []
    }

    DHTFindPeersModel.prototype.onSubmit = function(event) {
        $.post('/api/peerSearch', {
            infohash: this.infohash
        });

        this.infohash = '';
        event.preventDefault();
    }

    return DHTFindPeersModel;
})

