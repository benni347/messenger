"use strict";

import { RetrieveEnvValues } from "../wailsjs/go/main/App.js";

let appId;
let appKey;
let clusterId;
let appSecret;

RetrieveEnvValues().then((env) => {
  console.log(`Env: ${env}`);
  appId = env.appId;
  appKey = env.appKey;
  clusterId = env.clusterId;
  appSecret = env.appSecret;
  console.info(`App Id: ${appId}`);
  console.info(`App Key: ${appKey}`);
  console.info(`App Cluster: ${clusterId}`);
});

console.info(`App Cluster: ${clusterId}`);
