(function() {
    'use strict';

    var base = 'http://127.0.0.1:47935/';

    var script = document.createElement('script');
    script.setAttribute(
        'src',
        '//ajax.googleapis.com/ajax/libs/jquery/1.11.0/jquery.min.js');
    script.addEventListener('load', function() {
        jQuery('head *, body *').remove();
        jQuery('head').append(
            $('<base />').attr({href: base}),
            $('<script />').attr({src: '/_s/js/client.js'})
        );
    })
    document.body.appendChild(script);
}());
