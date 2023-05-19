"use strict";

import { RetrieveEnvValues } from "../wailsjs/go/main/App.js";

let appId = "";
let appKey = "";
let clusterId = "";
let appSecret = "";
let channelName = "1";
const messageLog = document.getElementById("message-log");

RetrieveEnvValues().then((env) => {
  console.log(`Env: ${env}`);
  appId = env.appId;
  appKey = env.appKey;
  clusterId = env.clusterId;
  appSecret = env.appSecret;
  console.info(`App Id: ${appId}`);
  console.info(`App Key: ${appKey}`);
  console.info(`App Cluster: ${clusterId}`);
  PusherClient();
});

// Enable pusher logging - don't include this in production
Pusher.logToConsole = true;

function PusherClient() {
  const pusher = new Pusher(appKey, {
    cluster: clusterId,
  });
  // TODO: Add authentication, see https://pusher.com/docs/channels/using_channels/authentication
  // TODO: ADD the ability to send messages to the channel
  // TODO: add Private channels
  // TODO: Store the msg in a database
  // TODO: Add ability to send messages to a specific user
  // The format from the server should be: {"message": "message", "time": "time"}
  const channel = pusher.subscribe(channelName);
  channel.bind("msg-recive", function (data) {
    console.info(`Pusher data: ${JSON.stringify(data)}`);
    let msgParagragh = document.createElement("p");
    let msgText = data.message;
    let timeMsg = data.time;
    msgParagragh.innerHTML = `${msgText}`;
    console.info(`Message: ${msgText}`);
    console.info(`Time: ${timeMsg}`);
    messageLog.appendChild(msgParagragh);
  });
}
