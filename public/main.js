$(function() {
  loadPlayer("player", function(player) {

    $("#volume").change(function() { player.setVolume($("#volume").val()); });

    var wantPaused = false;

    player.addEventListener("onStateChange", function(s) {
      if (s.data == 2) {
        wantPaused = true;
      }
      if (s.data == 1) {
        wantPaused = false;
      }
    });

    function annimAvatar(id) {
      setTimeout(function() {
        $("#" + id)
            .css({
              left: (Math.random() * 80 + 10) + "%",
              transform: 'rotate(' + (Math.random() * 50 - 25) + 'deg)',
            })
            .animate({
              bottom: -50,
            })
            .delay(Math.random() * 1000)
            .animate(
                {
                  bottom: -250,
                },
                function() { annimAvatar(id); });
      }, Math.random() * 15000);
    }
    annimAvatar("mouki");
    annimAvatar("marsou");
    annimAvatar("rico");

    var PC = $('#playlist');
    var userName = localStorage.getItem("username");

    if (!userName) {
      $("#add-1").show();
      $("#add-2").hide();
    } else {
      $("#add-2").show();
      $("#add-1").hide();
    }

    $("#name-ok")
        .click(function() {
          userName = $("#name").val();
          if (!userName) {
            return;
          }

          localStorage.setItem("username", userName);
          $("#add-2").show();
          $("#add-1").hide();
        });

    $('#addb')
        .click(function() {
          var id = parseYTU($('#addv').val());
          if (!id) {
            return;
          }
          $.post('/add?id=' + id + '&user=' + userName).fail(function(d) {
            alert(d.responseText);
          });
          $('#addv').val('');
        });

    function buildPL(pl) {
      PC.empty();

      if (pl.items) {
        for (var i in pl.items) {
          PC.append(buildPlItem(pl.items[i]));
        }
      }
    }

    function buildPlItem(item) {
      var e = $(
          '<li>' +
          '<div class="title"><b>' + item.video.title + '</b></div>' +
          '<div class="added-by">Ajouté par ' + item.added_by.name + '</div>' +
          '<div class="likes"><b>' + item.likes.length +
          ' votes!</b> <button class="vote">JE VOTE !</button></div>' +
          '<img src="sep.png">' +
          '</li>');

      e.find(".vote").click(function() {
        if (!userName) {
          alert("Met ton nom pour voter ! #GDPR");
          return;
        }
        $.post("/like?id=" + item.video.id + "&user=" + userName)
            .fail(function(d) { alert(d.responseText); });
      });

      return e;
    }

    onPlaylist(function(pl) {
      buildPL(pl);
      if (!pl.items || !pl.items.length) {
        return;
      }
      var plCurrent = pl.items[0];
      var pCurrentId = parseYTU(player.getVideoUrl());

      if (plCurrent.video.id != pCurrentId) {
        var start =
            (new Date().getTime() - Date.parse(pl.current_started_at)) / 1000;

        player[wantPaused ? "cueVideoById" : "loadVideoById"](
            plCurrent.video.id, start);
      }
    });
  });
});

function parseYTU(u) {
  var f = u.match(/[?&]v=([^&]+)/);
  if (!f) {
    return null;
  }

  return f[1];
}
