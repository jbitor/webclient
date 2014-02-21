'use strict';

angular.module('jbitor.client', [])

.controller('jbitorClientController', function(
    $scope,
    makeDHTMonitor
) {
    // TODO(JB): Stop passing $scope.
    $scope.dht = makeDHTMonitor($scope);
})

.factory('makeDHTMonitor', function() {
    return function makeDHTMonitor($scope) {
        var monitor = {
            nodeCounts: {
                GoodNodes: 0,
                UnknownNodes: 0,
                BadNodes: 0,
            },
            health: function() {
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
        };

        setInterval(updateNodeCounts, 1000);

        function updateNodeCounts() {
            $.getJSON('/api/nodeCounts.json').then(function(nodeCounts) {
                $scope.$apply(function() {
                    monitor.nodeCounts = nodeCounts;
                });
            }, function(err) {
                console.error("Failed to get update node counts", err);
            });
        }

        return monitor;
    };
});
