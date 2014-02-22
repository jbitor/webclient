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
    $scope.dht.startUpdating(750, $scope);

    $scope.findPeers = new DHTFindPeersModel();
})

.factory('DHTFindPeersModel', function() {
    function DHTFindPeersModel() {
        this.infohash = '';

        this.requests = []
    }

    DHTFindPeersModel.prototype.onSubmit = function(event) {
        $.post('/api/peerRequest', {
            infohash: this.infohash
        });

        this.infohash = '';
        event.preventDefault();
    }

    return DHTFindPeersModel;
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

        $.getJSON('/api/clientState.json').then(function(state) {
            $scope.$apply(function() {
                that.nodeCounts = state['nodeCounts'];
                $scope.findPeers.requests = state['peerRequests'];
            });
        }, function(err) {
            console.error("Failed to get update node counts", err);
        });
    }

    return DHTStatusModel;
});
