function onPlaylist(cb) {
  conn = new WebSocket('ws://' + document.location.host + '/ws');
  conn.onclose = function(evt) { console.warn('connection closed'); };
  conn.onmessage = function(evt) {
    console.log('got message :', evt.data);
    cb(JSON.parse(evt.data));
  };
}