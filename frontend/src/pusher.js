"use strict";

import { RetrieveEnvValues } from "../wailsjs/go/main/App.js";

let appId;
let appKey;
let clusterId;
let appSecret;
let channelName = "my-channel";
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

  const channel = pusher.subscribe(channelName);
  channel.bind("my-event", function(data) {
    console.info(`Pusher data: ${JSON.stringify(data)}`);
    messageLog.append(`${JSON.stringify(data)}\n`);
    alert(JSON.stringify(data));
  });
}
