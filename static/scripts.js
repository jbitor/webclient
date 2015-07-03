document.addEventListener('DOMContentLoaded', function() {
  'use strict';
  var input = document.querySelector('input');
  input.focus();
  input.setSelectionRange(input.value.length, input.value.length);

  setInterval(function() {
    fetch('/api/clientState.json').then(function(response) {
      return response.json();
    }).then(function(responseData) {
      document.querySelector('#clientState').textContent =
          JSON.stringify(responseData, null, 2);
    });
  }, 4000);

  document.querySelector('form').addEventListener('submit', function(event) {
    event.preventDefault();
    var data = new FormData;
    data.append('infohash', document.querySelector('input').value);
    fetch('/api/peerSearch', {
      method: 'POST',
      body: data
    });
    document.querySelector('input').value = '';
  })
});