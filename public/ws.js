function onPlaylist(cb) {
  function connect() {
    console.log("connecting");
    conn = new WebSocket('ws://' + document.location.host + '/ws');
    conn.onopen = function(event) { console.log("connected"); };
    conn.onclose = function(evt) {
      console.warn('connection closed', arguments);
      setTimeout(connect, 2000);
    };
    conn.onmessage = function(evt) {
      console.log('got message :', evt.data);
      cb(JSON.parse(evt.data));
    };
  }
  connect();
}