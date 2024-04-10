conn = new WebSocket("ws://" + document.location.host + "/ws");
conn.onerror = console.error;
conn.onopen = (e) => console.log("Connected to dev server", e);
conn.onclose = function (evt) {
  console.log("Connection Closed");
  setTimeout(function () {
    location.reload();
  }, 2000);
};
