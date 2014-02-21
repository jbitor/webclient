'use strict';

angular.module('jbitor.client', [])

.controller('jbitorClientController', function(
    $scope
) {
    $scope.dht = {
        nodeCounts: {}
    };

    setInterval(updateNodeCounts, 1000);

    function updateNodeCounts() {
        console.log("Requesting updated node counts.");

        $.getJSON('/api/nodeCounts.json').then(function(nodeCounts) {
            console.log("Got updated node counts.", nodeCounts);

            $scope.$apply(function() {
                $scope.dht.nodeCounts = nodeCounts;
            });
        }, function(err) {
            console.error("Failed to get update node counts", err);
        });
    }
});
