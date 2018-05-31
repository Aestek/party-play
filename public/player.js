function loadPlayer(id, cb) {
  var tag = document.createElement('script');

  tag.src = "https://www.youtube.com/iframe_api";
  var firstScriptTag = document.getElementsByTagName('script')[0];
  firstScriptTag.parentNode.insertBefore(tag, firstScriptTag);

  window.onYouTubeIframeAPIReady = function() {
    var player = new YT.Player('player', {
      height: '390',
      width: '640',
      videoId: '5qm8PH4xAss',
      playerVars: {
        controls: '0',
        cc_load_policy: '0',
        disablekb: '1',
        iv_load_policy: '3',
        modestbranding: '1',
      },
      events: {
        'onReady': function() {
          player.setVolume(80);
          cb(player);
        },
      }
    });
  };
}