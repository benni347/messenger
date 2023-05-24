"use strict";

/**
 * Retrieves the current chat room id from the 'body' element's 'data-current-chat-room-id' attribute.
 *
 * @returns {string} The current chat room id.
 */
function getChatRoomId() {
  return document.getElementById("body").attributes["data-current-chat-room-id"]
    .value;
}

/**
 * Validates the provided chat room id. Currently, only "00000000001" is considered valid.
 *
 * @param {string} chatRoomId - The id of the chat room to validate.
 * @returns {boolean} True if the chat room id is valid, false otherwise.
 */
function validateChatRoomId(chatRoomId) {
  if (chatRoomId === "00000000001") {
    return true;
  }
  return false;
}

/**
 * Adds a note to the 'person' div every 10 seconds. The note reminds users that the chat room is public, and that messages are stored unencrypted.
 */
function addNote() {
  if (!validateChatRoomId(getChatRoomId())) {
    return;
  }
  const noteP = document.createElement("p");
  noteP.innerHTML =
    "Note: This is a public chat room. Anyone can see your messages. The messages are stored in a database, unencrypted.";
  const personDiv = document.querySelector(".chat-note");

  personDiv.appendChild(noteP);
}

addNote();
