document.addEventListener('DOMContentLoaded', function() {
  'use strict';

  if (location.pathname == '/' && /^\?btih=/.test(location.search)) {
    window.location = '/' + location.search.slice(6);
    return;
  }

  var infohash = location.pathname.slice(1);

  if (!infohash) {
    return;
  }

  var input = document.querySelector('input');
  input.focus();
  input.value = infohash;
  input.setSelectionRange(0, input.value.length);

  var magnetLink = document.createElement('a');
  magnetLink.textContent = 'Magnet Link';
  magnetLink.href = 'magnet:?xt=urn:btih:' + infohash;
  document.querySelector('.nav').appendChild(magnetLink);

  var torrent = document.createElement('a');
  torrent.textContent = 'Download Torrent';
  torrent.href = '/' + infohash + '.torrent';
  document.querySelector('.nav').appendChild(torrent);

  var info = document.querySelector('#info');
  var infoPre = document.createElement('pre');
  infoPre.textContent = "Loading torrent info...";
  info.appendChild(infoPre);

  fetch('/' + infohash + '.json').then(function(response) {
    return response.json();
  }).then(function(data) {
    infoPre.textContent = JSON.stringify(data, null, 2);
  })
});