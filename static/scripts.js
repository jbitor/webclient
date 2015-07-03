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

  function peerSearch(infohash) {
    var data = new FormData;
    data.append('infohash', infohash);
    return fetch('/api/peerSearch', {
      method: 'POST',
      body: data
    }).then(function(response) {
      console.info("Started peer search.");
    }, function(error) {
      console.error("Unable to initiate peer search due to error:", error);
      throw error;
    });
  }

  window.peerSearch = peerSearch;
  console.info(
      "Use %cpeerSearch(infohash)%c to start a search.", 'color: blue;', 'color: auto;');

  document.querySelector('form').addEventListener('submit', function(event) {
    event.preventDefault();
    var input = document.querySelector('input');
    input.disabled = true;
    input.classList.remove('successful', 'error')
    peerSearch(input.value).then(function(){
      input.classList.add('successful');
      input.disabled = false;
      input.setSelectionRange(0, input.value.length);
    }, function() {
      input.classList.add('error');
      input.disabled = false;
      input.setSelectionRange(0, input.value.length);
    });
  });
});