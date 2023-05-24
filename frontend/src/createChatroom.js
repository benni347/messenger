"use strict";

import { CreateChatRoomId } from "../wailsjs/go/main/App.js";

function openNewChatRoomDialog() {
  const newChatRoomDialog = document.getElementById("new_chat_room_dialog");
  newChatRoomDialog.showModal();
}

window.addEventListener("DOMContentLoaded", () => {
  const newChatRoomBtn = document.getElementById("new_chat_room_wrapper");
  if (newChatRoomBtn) {
    newChatRoomBtn.addEventListener("click", () => {
      CreateChatRoomId();
    });
  }
});
