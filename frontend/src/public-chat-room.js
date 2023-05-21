"use strict";

function getChatRoomId() {
  return document.getElementById("body").attributes["data-current-chat-room-id"]
    .value;
}

function validateChatRoomId(chatRoomId) {
  if (chatRoomId === "00000000001") {
    return true;
  }
  return false;
}

function addNote() {
  const noteP = document.createElement("p");
  noteP.innerHTML =
    "Note: This is a public chat room. Anyone can see your messages. The messages are stored in a database, unencrypted.";
  noteP.style.gridArea = "notes";
  const personDiv = document.querySelector(".person");

  personDiv.appendChild(noteP);
}

setInterval(addNote(), 10 * 1000);
