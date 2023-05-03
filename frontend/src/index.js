"use strict";

import { Connect, CreatePeerConnection Dissconnect, ReceiveLocalMessage, ReceiveRemoteMessage, Send } from './wailsjs/go/main/App.js';

class Messaging {
  constructor() {
    super();
    this.localMessage = "";
    this.remoteMessage = "";
  }
}
