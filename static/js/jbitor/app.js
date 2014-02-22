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
    DHTStatusModel,
    DHTFindPeersModel
) {
    $scope.dht = new DHTStatusModel();
    $scope.dht.startUpdating(2500, $scope);

    $scope.findPeers = new DHTFindPeersModel();
})

.factory('DHTFindPeersModel', function(
    DHTFindPeersRequest
) {
    function DHTFindPeersModel() {
        this.infohash = '';

        this.requests = [
            new DHTFindPeersRequest('3c44dd30710c4d98d8ded1612428d7f9b3a6e44e'),
            new DHTFindPeersRequest('612428d7f9b3a6e44e4dd30710c4d98d8ded13c4'),
            new DHTFindPeersRequest('e4dd30710c4d98d8ded13c4612428d7f9b3a6e44')
        ]

        this.requests[1].peers = [7];
        this.requests[2].peers = [1, 2, 3];
    }

    DHTFindPeersModel.prototype.onSubmit = function(event) {
        this.infohash = '';
        event.preventDefault();
    }

    return DHTFindPeersModel;
})

.factory('DHTFindPeersRequest', function() {
    function DHTFindPeersRequest(infohash) {
        this.infohash = infohash;
        this.peers = [];
    }

    return DHTFindPeersRequest
})

.factory('DHTStatusModel', function(
    $interval
) {
    function DHTStatusModel(updateInterval) {
        this.nodeCounts = {
            GoodNodes: 0,
            UnknownNodes: 0,
            BadNodes: 0,
        };
    }

    DHTStatusModel.prototype.health = function() {
        if (this.nodeCounts.GoodNodes) {
            return Math.min(
                1.0,
                Math.max(
                    0.1,
                    this.nodeCounts.GoodNodes / 32
                )
            );
        } else if (this.nodeCounts.UnknownNodes) {
            return 0.05;
        }
    }

    DHTStatusModel.prototype.startUpdating = function(interval, $scope) {
        var that = this;

        var intervalPromise = $interval(function() {
            that.updateNodeCounts($scope);
        }, interval);

        $scope.$on('$destroy', function() {
            $interval.cancel(intervalPromise);
        });
    }

    DHTStatusModel.prototype.updateNodeCounts = function($scope) {
        var that = this;

        $.getJSON('/api/nodeCounts.json').then(function(nodeCounts) {
            $scope.$apply(function() {
                that.nodeCounts = nodeCounts;
            });
        }, function(err) {
            console.error("Failed to get update node counts", err);
        });
    }

    return DHTStatusModel;
});
