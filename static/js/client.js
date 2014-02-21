'use strict';

angular.module('jbitor.client', [])

.controller('jbitorClientController', function(
    $scope,
    DHTModel,
    FindPeersModel
) {
    $scope.dht = new DHTModel();
    $scope.dht.startUpdating(2500, $scope);

    $scope.findPeers = new FindPeersModel();
})

.factory('FindPeersModel', function() {
    function FindPeersModel() {
        this.infohash = '';
    }

    FindPeersModel.prototype.onSubmit = function(event) {
        this.infohash = '';
        event.preventDefault();
    }

    return FindPeersModel;
})

.factory('DHTModel', function(
    $interval
) {
    function DHTModel(updateInterval) {
        this.nodeCounts = {
            GoodNodes: 0,
            UnknownNodes: 0,
            BadNodes: 0,
        };
    }

    DHTModel.prototype.health = function() {
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

    DHTModel.prototype.startUpdating = function(interval, $scope) {
        var that = this;

        var intervalPromise = $interval(function() {
            that.updateNodeCounts($scope);
        }, interval);

        $scope.$on('$destroy', function() {
            $interval.cancel(intervalPromise);
        });
    }

    DHTModel.prototype.updateNodeCounts = function($scope) {
        var that = this;

        $.getJSON('/api/nodeCounts.json').then(function(nodeCounts) {
            $scope.$apply(function() {
                that.nodeCounts = nodeCounts;
            });
        }, function(err) {
            console.error("Failed to get update node counts", err);
        });
    }

    return DHTModel;
});
