(function() {
    'use strict';

    jQuery;
    $;

    var base = 'http://127.0.0.1:47935/';

    $('head').append(
        $('<link>').attr({
            rel: 'stylesheet',
            href: base + '/_s/css/style.css'})
    );

    var goodNodesEl, badNodesEl, unknownNodesEl;

    $('body').append(
        $('<section>').append(
            $('<h1>').text('jbitor'),

            $('<ul>').append(
                $('<li>').append(
                    "Good nodes: ",
                    goodNodesEl = $('<span>-</span>')
                ),
                $('<li>').append(
                    "Unknown nodes: ",
                    unknownNodesEl = $('<span>-</span>')
                ),
                $('<li>').append(
                    "Bad nodes: ",
                    badNodesEl = $('<span>-</span>')
                )
            )
        )
    );

    setInterval(updateNodeCounts, 1000);

    function updateNodeCounts() {
        $.get(base + 'api/nodeCounts.json').then(function(data) {
            goodNodesEl.text(data.GoodNodes);
            unknownNodesEl.text(data.UnknownNodes);
            badNodesEl.text(data.BadNodes);
        }, function(err) {
            $('body section ul li span').text(
                'error: ' + JSON.stringify(err));
        });
    }
}());
