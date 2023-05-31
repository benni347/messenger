"use strict";

import { RetrieveEnvValues } from "../wailsjs/go/main/App.js";
import Pusher from "pusher-js";

const channelName = "1";
// Retrieve the input and button elements
const messageInput = document.getElementById("message-input");
const submitButton = document.getElementById("submit");

RetrieveEnvValues().then((env) => {
  const appKey = env.appKey;
  const clusterId = env.clusterId;
  PusherClient(appKey, clusterId);
});

/**
 * Initializes a new Pusher client and sets up a subscription to a specified channel.
 *
 * This function first creates a new instance of the Pusher class with the provided 'appKey'
 * and 'clusterId'. Then, it subscribes to a channel specified by 'channelName'. When a
 * 'msg-recive' event is triggered on the subscribed channel, the function logs the received
 * data to the console, creates a new 'p' HTML element, populates it with the message text,
 * logs the message and time to the console, and appends the new element to the 'messageLog'
 * element.
 *
 * Note:
 * 1. The HTML elements used in this function must exist in the HTML document before this
 *    function is called.
 * 2. The 'appKey', 'clusterId', 'channelName', and 'messageLog' are external dependencies
 *    to this function and must be correctly set.
 * 3. This function is a starting point and currently has several TODO items, including
 *    implementing authentication, the ability to send messages, the use of private channels,
 *    storing messages in a database, and the ability to send messages to specific users.
 *
 * @function PusherClient
 */

function PusherClient(appKey, clusterId) {
  const pusher = new Pusher(appKey, {
    cluster: clusterId,
  });

  // The format from the server should be: {"message": "message", "time": "time"}
  const channel = pusher.subscribe(channelName);
  channel.bind("msg-recive", (data) => {
    console.info(`Pusher data: ${JSON.stringify(data)}`);
    const msgParagragh = document.createElement("p");
    const msgText = data.message;
    const timeMsg = data.time;
    msgParagragh.innerHTML = `${msgText}`;
    console.info(`Message: ${msgText}`);
    console.info(`Time: ${timeMsg}`);
    messageLog.appendChild(msgParagragh);
  });

  // Add event listener to the input field to listen for the 'Enter' key
  messageInput.addEventListener("keyup", function (event) {
    // Check if the 'Enter' key was pressed
    if (event.key === "Enter") {
      sendMessage(channel);
    }
  });

  // Add event listener to the button to listen for clicks
  submitButton.addEventListener("click", function () {
    sendMessage(channel);
  });
}

// Function to send message
function sendMessage(channel) {
  // Get the message from the input field
  const message = messageInput.value;

  // Check if the message is not empty
  if (message.trim() !== "") {
    // Trigger a 'client-msg-send' event on the channel with the message as the data
    channel.trigger("client-msg-send", { message: message });

    // Clear the input field
    messageInput.value = "";
  }
}
