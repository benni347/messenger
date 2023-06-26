"use strict";

const socket = new WebSocket("ws://s.rabbitmq.cedric.5y5.one/websocket");

socket.onmessage = function (event) {
  const otherMessage = event.data;
  const messageLog = document.getElementById("message-log");
  const messageDiv = document.createElement("div");
  const messageUsernameDiv = document.createElement("div");
  const messageTextDiv = document.createElement("div");
  const otherUsername = "other";
  messageUsernameDiv.innerHTML = otherUsername;
  messageTextDiv.innerHTML = otherMessage;
  messageTextDiv.className = "text";
  messageUsernameDiv.className = "username";
  messageDiv.className = "message";
  messageDiv.appendChild(messageUsernameDiv);
  messageDiv.appendChild(messageTextDiv);
  messageLog.appendChild(messageDiv);
  messageLog.scrollTop = messageLog.scrollHeight;
};
