(() => {
  (function () {
    let n = function (e) {
        for (var s = window[e.shift()]; s && e.length; ) s = s[e.shift()];
        return s;
      },
      o = n(["chrome", "webview", "postMessage"]),
      t = n(["webkit", "messageHandlers", "external", "postMessage"]);
    if (!o && !t) {
      console.error("Unsupported Platform");
      return;
    }
    o && (window.WailsInvoke = (e) => window.chrome.webview.postMessage(e)),
      t &&
        (window.WailsInvoke = (e) =>
          window.webkit.messageHandlers.external.postMessage(e));
  })();
})();
