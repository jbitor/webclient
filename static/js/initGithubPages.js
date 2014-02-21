(function() {
    'use strict';

    var iframe = document.createElement('iframe');
    iframe.src = 'http://127.0.0.1:47935/';
    iframe.style.position = 'absolute';
    iframe.style.top = '0';
    iframe.style.bottom = '0';
    iframe.style.height = '100%';
    iframe.style.left = '0';
    iframe.style.right = '0';
    iframe.style.width = '100%';
    iframe.style.zIndex = 800;

    document.body.appendChild(iframe);
}());
