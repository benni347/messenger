(() => {
  var L = Object.defineProperty;
  var c = (e, n) => {
    for (var o in n) L(e, o, { get: n[o], enumerable: !0 });
  };
  var m = {};
  c(m, {
    LogDebug: () => M,
    LogError: () => A,
    LogFatal: () => H,
    LogInfo: () => P,
    LogLevel: () => J,
    LogPrint: () => R,
    LogTrace: () => z,
    LogWarning: () => B,
    SetLogLevel: () => G,
  });
  function d(e, n) {
    window.WailsInvoke("L" + e + n);
  }
  function z(e) {
    d("T", e);
  }
  function R(e) {
    d("P", e);
  }
  function M(e) {
    d("D", e);
  }
  function P(e) {
    d("I", e);
  }
  function B(e) {
    d("W", e);
  }
  function A(e) {
    d("E", e);
  }
  function H(e) {
    d("F", e);
  }
  function G(e) {
    d("S", e);
  }
  var J = { TRACE: 1, DEBUG: 2, INFO: 3, WARNING: 4, ERROR: 5 };
  var v = class {
      constructor(n, o, i) {
        (this.eventName = n),
          (this.maxCallbacks = i || -1),
          (this.Callback = (t) => (
            o.apply(null, t),
            this.maxCallbacks === -1
              ? !1
              : ((this.maxCallbacks -= 1), this.maxCallbacks === 0)
          ));
      }
    },
    l = {};
  function p(e, n, o) {
    l[e] = l[e] || [];
    let i = new v(e, n, o);
    return l[e].push(i), () => F(i);
  }
  function b(e, n) {
    return p(e, n, -1);
  }
  function E(e, n) {
    return p(e, n, 1);
  }
  function S(e) {
    let n = e.name;
    if (l[n]) {
      let o = l[n].slice();
      for (let i = 0; i < l[n].length; i += 1) {
        let t = l[n][i],
          r = e.data;
        t.Callback(r) && o.splice(i, 1);
      }
      o.length === 0 ? g(n) : (l[n] = o);
    }
  }
  function y(e) {
    let n;
    try {
      n = JSON.parse(e);
    } catch {
      let i = "Invalid JSON passed to Notify: " + e;
      throw new Error(i);
    }
    S(n);
  }
  function C(e) {
    let n = { name: e, data: [].slice.apply(arguments).slice(1) };
    S(n), window.WailsInvoke("EE" + JSON.stringify(n));
  }
  function g(e) {
    delete l[e], window.WailsInvoke("EX" + e);
  }
  function D(e, ...n) {
    g(e),
      n.length > 0 &&
        n.forEach((o) => {
          g(o);
        });
  }
  function F(e) {
    let n = e.eventName;
    (l[n] = l[n].filter((o) => o !== e)), l[n].length === 0 && g(n);
  }
  var f = {};
  function U() {
    var e = new Uint32Array(1);
    return window.crypto.getRandomValues(e)[0];
  }
  function j() {
    return Math.random() * 9007199254740991;
  }
  var W;
  window.crypto ? (W = U) : (W = j);
  function s(e, n, o) {
    return (
      o == null && (o = 0),
      new Promise(function (i, t) {
        var r;
        do r = e + "-" + W();
        while (f[r]);
        var a;
        o > 0 &&
          (a = setTimeout(function () {
            t(Error("Call to " + e + " timed out. Request ID: " + r));
          }, o)),
          (f[r] = { timeoutHandle: a, reject: t, resolve: i });
        try {
          let u = { name: e, args: n, callbackID: r };
          window.WailsInvoke("C" + JSON.stringify(u));
        } catch (u) {
          console.error(u);
        }
      })
    );
  }
  window.ObfuscatedCall = (e, n, o) => (
    o == null && (o = 0),
    new Promise(function (i, t) {
      var r;
      do r = e + "-" + W();
      while (f[r]);
      var a;
      o > 0 &&
        (a = setTimeout(function () {
          t(Error("Call to method " + e + " timed out. Request ID: " + r));
        }, o)),
        (f[r] = { timeoutHandle: a, reject: t, resolve: i });
      try {
        let u = { id: e, args: n, callbackID: r };
        window.WailsInvoke("c" + JSON.stringify(u));
      } catch (u) {
        console.error(u);
      }
    })
  );
  function T(e) {
    let n;
    try {
      n = JSON.parse(e);
    } catch (t) {
      let r = `Invalid JSON passed to callback: ${t.message}. Message: ${e}`;
      throw (runtime.LogDebug(r), new Error(r));
    }
    let o = n.callbackid,
      i = f[o];
    if (!i) {
      let t = `Callback '${o}' not registered!!!`;
      throw (console.error(t), new Error(t));
    }
    clearTimeout(i.timeoutHandle),
      delete f[o],
      n.error ? i.reject(n.error) : i.resolve(n.result);
  }
  window.go = {};
  function O(e) {
    try {
      e = JSON.parse(e);
    } catch (n) {
      console.error(n);
    }
    (window.go = window.go || {}),
      Object.keys(e).forEach((n) => {
        (window.go[n] = window.go[n] || {}),
          Object.keys(e[n]).forEach((o) => {
            (window.go[n][o] = window.go[n][o] || {}),
              Object.keys(e[n][o]).forEach((i) => {
                window.go[n][o][i] = (function () {
                  let t = 0;
                  function r() {
                    let a = [].slice.call(arguments);
                    return s([n, o, i].join("."), a, t);
                  }
                  return (
                    (r.setTimeout = function (a) {
                      t = a;
                    }),
                    (r.getTimeout = function () {
                      return t;
                    }),
                    r
                  );
                })();
              });
          });
      });
  }
  var x = {};
  c(x, {
    WindowCenter: () => q,
    WindowFullscreen: () => Z,
    WindowGetPosition: () => se,
    WindowGetSize: () => ne,
    WindowHide: () => le,
    WindowIsFullscreen: () => _,
    WindowIsMaximised: () => ue,
    WindowIsMinimised: () => pe,
    WindowIsNormal: () => We,
    WindowMaximise: () => we,
    WindowMinimise: () => ce,
    WindowReload: () => N,
    WindowReloadApp: () => V,
    WindowSetAlwaysOnTop: () => te,
    WindowSetBackgroundColour: () => me,
    WindowSetDarkTheme: () => $,
    WindowSetLightTheme: () => Y,
    WindowSetMaxSize: () => oe,
    WindowSetMinSize: () => ie,
    WindowSetPosition: () => re,
    WindowSetSize: () => ee,
    WindowSetSystemDefaultTheme: () => X,
    WindowSetTitle: () => Q,
    WindowShow: () => ae,
    WindowToggleMaximise: () => de,
    WindowUnfullscreen: () => K,
    WindowUnmaximise: () => fe,
    WindowUnminimise: () => ge,
  });
  function N() {
    window.location.reload();
  }
  function V() {
    window.WailsInvoke("WR");
  }
  function X() {
    window.WailsInvoke("WASDT");
  }
  function Y() {
    window.WailsInvoke("WALT");
  }
  function $() {
    window.WailsInvoke("WADT");
  }
  function q() {
    window.WailsInvoke("Wc");
  }
  function Q(e) {
    window.WailsInvoke("WT" + e);
  }
  function Z() {
    window.WailsInvoke("WF");
  }
  function K() {
    window.WailsInvoke("Wf");
  }
  function _() {
    return s(":wails:WindowIsFullscreen");
  }
  function ee(e, n) {
    window.WailsInvoke("Ws:" + e + ":" + n);
  }
  function ne() {
    return s(":wails:WindowGetSize");
  }
  function oe(e, n) {
    window.WailsInvoke("WZ:" + e + ":" + n);
  }
  function ie(e, n) {
    window.WailsInvoke("Wz:" + e + ":" + n);
  }
  function te(e) {
    window.WailsInvoke("WATP:" + (e ? "1" : "0"));
  }
  function re(e, n) {
    window.WailsInvoke("Wp:" + e + ":" + n);
  }
  function se() {
    return s(":wails:WindowGetPos");
  }
  function le() {
    window.WailsInvoke("WH");
  }
  function ae() {
    window.WailsInvoke("WS");
  }
  function we() {
    window.WailsInvoke("WM");
  }
  function de() {
    window.WailsInvoke("Wt");
  }
  function fe() {
    window.WailsInvoke("WU");
  }
  function ue() {
    return s(":wails:WindowIsMaximised");
  }
  function ce() {
    window.WailsInvoke("Wm");
  }
  function ge() {
    window.WailsInvoke("Wu");
  }
  function pe() {
    return s(":wails:WindowIsMinimised");
  }
  function We() {
    return s(":wails:WindowIsNormal");
  }
  function me(e, n, o, i) {
    let t = JSON.stringify({ r: e || 0, g: n || 0, b: o || 0, a: i || 255 });
    window.WailsInvoke("Wr:" + t);
  }
  var k = {};
  c(k, { ScreenGetAll: () => ve });
  function ve() {
    return s(":wails:ScreenGetAll");
  }
  var h = {};
  c(h, { BrowserOpenURL: () => xe });
  function xe(e) {
    window.WailsInvoke("BO:" + e);
  }
  var I = {};
  c(I, { ClipboardGetText: () => he, ClipboardSetText: () => ke });
  function ke(e) {
    return s(":wails:ClipboardSetText", [e]);
  }
  function he() {
    return s(":wails:ClipboardGetText");
  }
  function Ie() {
    window.WailsInvoke("Q");
  }
  function be() {
    window.WailsInvoke("S");
  }
  function Ee() {
    window.WailsInvoke("H");
  }
  function Se() {
    return s(":wails:Environment");
  }
  window.runtime = {
    ...m,
    ...x,
    ...h,
    ...k,
    ...I,
    EventsOn: b,
    EventsOnce: E,
    EventsOnMultiple: p,
    EventsEmit: C,
    EventsOff: D,
    Environment: Se,
    Show: be,
    Hide: Ee,
    Quit: Ie,
  };
  window.wails = {
    Callback: T,
    EventsNotify: y,
    SetBindings: O,
    eventListeners: l,
    callbacks: f,
    flags: {
      disableScrollbarDrag: !1,
      disableWailsDefaultContextMenu: !1,
      enableResize: !1,
      defaultCursor: null,
      borderThickness: 6,
      shouldDrag: !1,
      deferDragToMouseMove: !1,
      cssDragProperty: "--wails-draggable",
      cssDragValue: "drag",
    },
  };
  window.wailsbindings &&
    (window.wails.SetBindings(window.wailsbindings),
    delete window.wails.SetBindings);
  delete window.wailsbindings;
  var ye = function (e) {
    var n = window
      .getComputedStyle(e.target)
      .getPropertyValue(window.wails.flags.cssDragProperty);
    return (
      n && (n = n.trim()),
      !(
        n !== window.wails.flags.cssDragValue ||
        e.buttons !== 1 ||
        e.detail !== 1
      )
    );
  };
  window.wails.setCSSDragProperties = function (e, n) {
    (window.wails.flags.cssDragProperty = e),
      (window.wails.flags.cssDragValue = n);
  };
  window.addEventListener("mousedown", (e) => {
    if (window.wails.flags.resizeEdge) {
      window.WailsInvoke("resize:" + window.wails.flags.resizeEdge),
        e.preventDefault();
      return;
    }
    if (ye(e)) {
      if (
        window.wails.flags.disableScrollbarDrag &&
        (e.offsetX > e.target.clientWidth || e.offsetY > e.target.clientHeight)
      )
        return;
      window.wails.flags.deferDragToMouseMove
        ? (window.wails.flags.shouldDrag = !0)
        : (e.preventDefault(), window.WailsInvoke("drag"));
      return;
    } else window.wails.flags.shouldDrag = !1;
  });
  window.addEventListener("mouseup", () => {
    window.wails.flags.shouldDrag = !1;
  });
  function w(e) {
    (document.documentElement.style.cursor =
      e || window.wails.flags.defaultCursor),
      (window.wails.flags.resizeEdge = e);
  }
  window.addEventListener("mousemove", function (e) {
    if (window.wails.flags.shouldDrag)
      if ((e.buttons !== void 0 ? e.buttons : e.which) <= 0)
        window.wails.flags.shouldDrag = !1;
      else {
        window.WailsInvoke("drag");
        return;
      }
    if (!window.wails.flags.enableResize) return;
    window.wails.flags.defaultCursor == null &&
      (window.wails.flags.defaultCursor =
        document.documentElement.style.cursor),
      window.outerWidth - e.clientX < window.wails.flags.borderThickness &&
        window.outerHeight - e.clientY < window.wails.flags.borderThickness &&
        (document.documentElement.style.cursor = "se-resize");
    let n = window.outerWidth - e.clientX < window.wails.flags.borderThickness,
      o = e.clientX < window.wails.flags.borderThickness,
      i = e.clientY < window.wails.flags.borderThickness,
      t = window.outerHeight - e.clientY < window.wails.flags.borderThickness;
    !o && !n && !i && !t && window.wails.flags.resizeEdge !== void 0
      ? w()
      : n && t
      ? w("se-resize")
      : o && t
      ? w("sw-resize")
      : o && i
      ? w("nw-resize")
      : i && n
      ? w("ne-resize")
      : o
      ? w("w-resize")
      : i
      ? w("n-resize")
      : t
      ? w("s-resize")
      : n && w("e-resize");
  });
  window.addEventListener("contextmenu", function (e) {
    window.wails.flags.disableWailsDefaultContextMenu && e.preventDefault();
  });
  window.WailsInvoke("runtime:ready");
})();
