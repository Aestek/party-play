$(function() {
  conn = new WebSocket('ws://' + document.location.host + '/ws');
  conn.onclose = function(evt) {
    console.warn('connection closed');
  };
  conn.onmessage = function(evt) {
    console.log('got message :', evt.data);
    buildPL(JSON.parse(evt.data));
  };

  var PC = $('#playlist');

  $('#addb').click(function() {
    var id = $('#addv').val();
    var name = $('#name').val();
    $.post('/add?id=' + id + '&user=' + name);
  });

  function buildPL(pl) {
    var plContent = '';

    if (pl.items) {
      for (var i in pl.items) {
        plContent += buildPlItem(pl.items[i]);
      }
    }

    PC.html(plContent);
  }

  function buildPlItem(item) {
    return '<li>' + item.video.title + '</li>'
  }
});