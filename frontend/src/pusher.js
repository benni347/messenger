"use strict";

import { RetrieveEnvValues } from "../wailsjs/go/main/App.js";

let appKey = "";
let clusterId = "";
const channelName = "1";
const messageLog = document.getElementById("message-log");

RetrieveEnvValues().then((env) => {
  appKey = env.appKey;
  clusterId = env.clusterId;
  PusherClient();
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
function PusherClient() {
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
}
