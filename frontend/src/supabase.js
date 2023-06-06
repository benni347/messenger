"use strict";

import {
  RetrieveEnvValues,
  ValidateEmail,
  GenerateUserName,
  CreateChatRoomId,
  GetOtherUserId,
} from "../wailsjs/go/main/App.js";

// Solved the fix me through importing it as a npm module
import { createClient } from "@supabase/supabase-js";

const options = {
  db: {
    schema: "public",
  },
  auth: {
    autoRefreshToken: true,
    persistSession: true,
    detectSessionInUrl: true,
  },
};

let supabaseKey = "";
let supabaseUrl = "";
let supabase = "";
let authenticated = false;

const messageLog = document.getElementById("message-log");
/**
 * Asynchronously signs in a user through email using Supabase authentication.
 *
 * This function first retrieves user input from the HTML elements with ids 'email-input',
 * 'password-input', and 'username'. It then attempts to sign the user in via Supabase's
 * auth.signIn method with the retrieved user input data. If an error occurs during the
 * signIn process, the error is logged to the console. Otherwise, the data received from
 * the signIn request is logged to the console.
 *
 * The function does not return anything.
 *
 * Note: The HTML elements used in this function must exist in the HTML document and be
 * populated with appropriate user input data (email, password, and username) before this
 * function is called.
 *
 * @async
 * @function signInThroughMail
 * @throws Will throw an error if the sign in process encounters any issues.
 */
async function signInThroughMail() {
  const email = document.getElementById("email-input").value;
  if (!ValidateEmail(email)) {
    console.error(`The email ${email} is not valid.`);
    return;
  }
  const password = document.getElementById("password-input").value;
  const user_name = document.getElementById("username").value;

  const { data, error } = await supabase.auth.signInWithPassword({
    email, // equals to email: email
    password, // equals to password: password
    options: {
      data: {
        user_name, // equals to user_name: user_name
      },
    },
  });
  authenticated = true;

  localStorage.setItem("authenticated", authenticated);
  changeButton();
  if (error) {
    console.error(`An error occured during the login: ${error}`);
  }

  location.reload();
}
/**
 * Gets the current username. If it doesn't exist, generates a new one.
 *
 * @async
 * @returns {Promise<string>} The username.
 */
const getUsername = async () => {
  const previousUsername = localStorage.getItem("username");
  const user = await supabase.auth.getUser();
  if (
    user &&
    user.data &&
    user.data.user &&
    user.data.user.user_metadata &&
    user.data.user.user_metadata.user_name
  ) {
    const username = user.data.user.user_metadata.user_name;
    return username;
  } else {
    const username = previousUsername || GenerateUserName(4);
    return username;
  }
};
/**
 * Sets the current username in local storage and updates the username field in the document.
 *
 * @async
 * @returns {Promise<string>} The username.
 */
const setUsername = async () => {
  const username = await getUsername();
  localStorage.setItem("username", username);
  const usernameParagraph = document.createElement("p");
  usernameParagraph.style.gridArea = "username";
  usernameParagraph.innerHTML = username;
  const userNameDiv = document.getElementById("username");
  if (userNameDiv.childElementCount > 0) {
    userNameDiv.removeChild(userNameDiv.childNodes[0]);
  }
  userNameDiv.appendChild(usernameParagraph);
  return username;
};

/**
 * Asynchronously signs out a user using Supabase authentication.
 *
 * This function attempts to sign the user out via Supabase's auth.signOut method. If an error
 * occurs during the signOut process, the error is logged to the console. The function does
 * not return anything.
 *
 * Note: A user must be signed in before this function is called.
 *
 * @async
 * @function signOut
 * @throws Will throw an error if the sign out process encounters any issues.
 */
async function signOut() {
  const { error } = await supabase.auth.signOut();
  if (error) {
    console.error(`An error occured during the logout: ${error}`);
  }
  authenticated = false;
  localStorage.setItem("authenticated", authenticated);
  localStorage.removeItem("username");
  setUsername();
  const body = document.getElementById("body");
  body.setAttribute("data-current-user-id", "");
  removeUserIdNote();
  changeButton();
}

/**
 * Asynchronously signs in a user using GitHub as OAuth provider via Supabase authentication.
 *
 * This function attempts to sign the user in via Supabase's auth.signInWithOAuth method with
 * 'github' as the OAuth provider. If an error occurs during the signIn process, the error
 * is logged to the console. Otherwise, the data received from the signIn request is logged
 * to the console.
 *
 * The function does not return anything.
 *
 * @async
 * @function signInTroughGithub
 * @throws Will throw an error if the sign in process encounters any issues.
 */
async function signInTroughGithub() {
  const { data, error } = await supabase.auth.signInWithOAuth({
    provider: "github",
  });
  if (error) {
    console.error(`An error occured during the login: ${error}`);
  }
  authenticated = true;
  localStorage.setItem("authenticated", authenticated);
  changeButton();
}

RetrieveEnvValues().then((env) => {
  supabaseKey = env.supaBaseApiKey;
  supabaseUrl = env.supaBaseUrl;
  supabase = createClient(supabaseUrl, supabaseKey, options);

  setUsername();
  addUserIdNote();
});

/**
 * Asynchronously registers a new user with Supabase authentication.
 *
 * This function first retrieves user input from the HTML elements with ids 'email-input',
 * 'password-input', and 'username'. It then attempts to sign the user up via Supabase's
 * auth.signUp method with the retrieved user input data. If an error occurs during the
 * signUp process, the error is logged to the console. Otherwise, the data received from
 * the signUp request is logged to the console.
 *
 * The function does not return anything.
 *
 * Note: The HTML elements used in this function must exist in the HTML document and be
 * populated with appropriate user input data (email, password, and username) before this
 * function is called.
 *
 * @async
 * @function signUp
 * @throws Will throw an error if the user creation process encounters any issues.
 */
async function signUp() {
  const email = document.getElementById("email-input-signup").value;
  const password = document.getElementById("password-signup").value;
  const user_name = document.getElementById("username").value;

  if (!ValidateEmail(email)) {
    console.error(`The email ${email} is not valid.`);
    return;
  }
  const { data, error } = await supabase.auth.signUp({
    email, // equals to email: email
    password, // equals to password: password
    options: {
      data: {
        user_name, // equals to user_name: user_name
      },
    },
  });
  if (error) {
    console.error(`An error occured during the creation of the user: ${error}`);
  }
}

/**
 * Asynchronously retrieves the authenticated user's ID from the Supabase database.
 *
 * If the user is authenticated according to the local storage, the function
 * attempts to retrieve the user object from Supabase. If the user object exists
 * and it contains the `id` field in its `data.user` property, this `id` is returned.
 * If the user is not authenticated or the user object doesn't have an `id` field,
 * the function returns `null`.
 *
 * @function
 * @async
 * @returns {Promise<string|null>} The authenticated user's ID as a string, or `null` if
 * the user is not authenticated or if the user object doesn't contain an `id`.
 * @example
 * const userId = await getId();
 * if (userId) {
 *   console.log("Authenticated user ID is", userId);
 * } else {
 *   console.log("User is not authenticated or doesn't have an ID.");
 * }
 */
const getId = async () => {
  if (
    localStorage.getItem("authenticated") === "true" ||
    localStorage.getItem("authenticated") === true
  ) {
    const user = await supabase.auth.getUser();
    if (user && user.data && user.data.user && user.data.user.id) {
      return user.data.user.id;
    } else {
      return null;
    }
  } else {
    return null;
  }
};

/**
 * Gets a range of messages from the 'messages' database table, ordered by timestamp in descending order.
 *
 * @async
 * @param {number} from - The starting index for the range of messages to retrieve.
 * @param {number} to - The ending index for the range of messages to retrieve.
 * @returns {Promise<Array>} The array of messages data.
 */
const getMessages = async (from, to) => {
  const { data } = await supabase
    .from("messages")
    .select()
    .range(from, to)
    .order("timestamp", { ascending: false });

  return data;
};

/**
 * Subscribes to the 'INSERT' event on the 'messages' database table and calls the provided handler function whenever a new message is inserted.
 *
 * @param {function} handler - The function to call when a new message is inserted.
 */
const onNewMessage = (handler) => {
  supabase
    .from("messages")
    .on("INSERT", (payload) => {
      handler(payload.new);
    })
    .subscribe();
};

/**
 * Inserts a new message into the 'messages' database table.
 *
 * @async
 * @param {string} username - The username of the sender of the message.
 * @param {string} text - The text of the message.
 * @returns {Promise<Object>} The data of the inserted message.
 */
const createNewMessage = async (username, text) => {
  const { data } = await supabase.from("messages").insert({ username, text });

  return data;
};

/**
 * Manages the messages in the chat, allowing to load older messages and send new ones.
 *
 * @returns {Object} An object containing the current username, a reference to the messages, a function to send messages, and a function to load older messages.
 */
const useMessages = () => {
  const username = getUsername();
  const messages = ref([]);
  const messagesCount = ref(0);
  const maxMessgesPerRequest = 50;
  /**
   * Loads a batch of messages from the server, updating the count of messages and appending the loaded messages to the existing ones.
   *
   * @async
   */
  const loadMessagesBatch = async () => {
    const loadedMessages = await getMessages(
      messagesCount.value,
      maxMessgesPerRequest - 1
    );

    messages.value = [...loadedMessages, ...messages.value];
    messagesCount.value += loadedMessages.length;
  };

  loadMessagesBatch();
  onNewMessage((newMessage) => {
    messages.value = [newMessage, ...messages.value];
    messagesCount.value += 1;
  });

  return {
    username,
    messages,
    async sendMessage(text) {
      if (text) {
        await createNewMessage(username, text);
      }
    },
    loadOlder() {
      loadMessagesBatch();
    },
  };
};

/**
 * Toggles the visibility of the sign-in and sign-out buttons.
 * The visibility is determined based on the 'authenticated' flag stored in local storage.
 * If the flag is 'true', the sign-in button is hidden and the sign-out button is displayed.
 * Otherwise, the sign-out button is hidden and the sign-in button is displayed.
 */
function changeButton() {
  if (
    localStorage.getItem("authenticated") === true ||
    localStorage.getItem("authenticated") === "true"
  ) {
    document.getElementById("signin-main-wrapper").style.display = "none";
    document.getElementById("signout-main-wrapper").style.display = "block";
  } else {
    document.getElementById("signin-main-wrapper").style.display = "block";
    document.getElementById("signout-main-wrapper").style.display = "none";
  }
}

/**
 * Displays the window for creating a new chat room.
 * This is achieved by removing the 'hidden' class from the 'new_chat_room_window' element.
 */
function newChatRoom() {
  const newChatRoomWindow = document.getElementById("new_chat_room_window");
  newChatRoomWindow.classList.remove("hidden");
}

/**
 * Creates a new chat room.
 * The function first gets the user id of the other person involved in the chat.
 * Then, it retrieves the current user's id by awaiting the result of the 'getId' function.
 * The retrieved ids (which may include dashes) are cleaned by removing all dashes.
 * A chat room id is then created using the cleaned ids, and this id is printed to the console.
 *
 * @async
 * @returns {Promise<void>} This function returns a promise that resolves to undefined.
 * It has no return value because the created chat room id is not returned, only logged.
 */
async function createNewChatRoom() {
  const other_user_id = document.getElementById("other_persons_uid").value;
  const myId = await getId(); // Await the promise returned by getId()

  const myIdWithoutDashes = myId.replace(/-/g, "");
  const otherIdWithoutDashes = other_user_id.replace(/-/g, "");

  const combindedIds = await CreateChatRoomId(
    myIdWithoutDashes,
    otherIdWithoutDashes
  );
  console.info(combindedIds);
  const body = document.querySelector("body");
  body.setAttribute("data-current-chat-room-id", combindedIds);
  localStorage.setItem("current-chat-room-id", combindedIds);
  addChatRoomId(combindedIds);
}

function addChatRoomId(newId) {
  // get existing ids
  let storedChatRoomIds = localStorage.getItem('all_chat_room_ids');
  let chatRoomIds = storedChatRoomIds ? JSON.parse(storedChatRoomIds) : [];

  // add new id
  chatRoomIds.push(newId);

  // store updated ids
  localStorage.setItem('all_chat_room_ids', JSON.stringify(chatRoomIds));
}

/**
 * Removes the last paragraph element from the '.note' div element.
 *
 * @function
 * @name removeUserIdNote
 * @example
 * // Removing the note
 * removeUserIdNote();
 */
function removeUserIdNote() {
  const noteDiv = document.querySelector(".note");
  const noteP = noteDiv.lastElementChild;

  if (noteP && noteP.tagName.toLowerCase() === "p") {
    noteDiv.removeChild(noteP);
  }
}

/**
 * Checks if the user is authenticated and if so, retrieves their user ID and adds a note on the page.
 * The authentication is based on the 'authenticated' flag stored in local storage.
 * If authenticated, the user ID is fetched using 'getId' function, which returns a Promise.
 * Once the Promise is resolved, the ID is added as a data attribute to the body element.
 * Furthermore, a paragraph element is created containing a user-friendly message which includes the user ID.
 * This paragraph is appended to the '.note' element on the page.
 */
function addUserIdNote() {
  if (
    localStorage.getItem("authenticated") === true ||
    localStorage.getItem("authenticated") === "true"
  ) {
    getId().then((id) => {
      const body = document.getElementById("body");
      body.setAttribute("data-current-user-id", id);
      const noteP = document.createElement("p");
      noteP.innerHTML =
        "Your User ID: " +
        id +
        "<br>Copy this ID and send it to your friend to start chatting!";
      const personDiv = document.querySelector(".note");

      personDiv.appendChild(noteP);
    });
  }
}

window.addEventListener("DOMContentLoaded", () => {
  changeButton();
  const signInWindow = document.getElementById("signin-window");
  const signUpWindow = document.getElementById("signup-window");
  const verifyEmailWindow = document.getElementById("verify-window");
  const verifyEmailBtn = document.getElementById("okay-verify-btn");
  const mainContentWrapper = document.getElementById("main-window-wrappper");
  const signOutBtn = document.getElementById("signout-main-wrapper");
  const signInBtn = document.getElementById("signing");
  const closeSignInBtn = document.getElementById("close-signin-button");
  const closeSignUpBtn = document.getElementById("close-signup-button");
  const closeNewChatRoomBtn = document.getElementById("close-new_chat-button");
  const signUpBtn = document.getElementById("signup-btn");
  const signUpBtnInSignInWindow = document.getElementById("signup");
  const signInBtnInSignUpWindow = document.getElementById("signin");
  const signInBtnInMainContentWrapper = document.getElementById(
    "signin-main-wrapper"
  );
  const githubBtn = document.getElementById("github-button");
  const githubBtnSignUp = document.getElementById("github-button-signup");
  const newChatRoomButtonOnMainContentWrapper = document.getElementById(
    "new_chat_room_wrapper"
  );
  const newChatRoomWindow = document.getElementById("new_chat_room_window");
  const newChatBtnInNewChatRoomWindow = document.getElementById("new_chat-btn");

  // Under this line define no more consts for html elements. Only function calls.
  if (signInBtn) {
    signInBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signInThroughMail();
    });
  }
  if (signOutBtn) {
    signOutBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signOut();
    });
  }
  if (signUpBtn) {
    signUpBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signUp();
      signUpWindow.classList.add("hidden");
      verifyEmailWindow.classList.remove("hidden");
    });
  }
  if (closeSignInBtn) {
    closeSignInBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signInWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
      signInWindow.classList.add("hidden");
    });
  }
  if (closeSignUpBtn) {
    closeSignUpBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signUpWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
    });
  }
  if (closeNewChatRoomBtn) {
    closeNewChatRoomBtn.addEventListener("click", (event) => {
      event.preventDefault();
      newChatRoomWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
    });
  }
  if (signUpBtnInSignInWindow) {
    signUpBtnInSignInWindow.addEventListener("click", (event) => {
      event.preventDefault();
      signInWindow.classList.add("hidden");
      signUpWindow.classList.remove("hidden");
    });
  }
  if (signInBtnInSignUpWindow) {
    signInBtnInSignUpWindow.addEventListener("click", (event) => {
      event.preventDefault();
      signUpWindow.classList.add("hidden");
      signInWindow.classList.remove("hidden");
    });
  }
  if (signInBtnInMainContentWrapper) {
    signInBtnInMainContentWrapper.addEventListener("click", (event) => {
      event.preventDefault();
      mainContentWrapper.classList.add("signin-portal");
      signInWindow.classList.remove("hidden");
    });
  }
  if (githubBtn) {
    githubBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signInTroughGithub();
    });
  }
  if (githubBtnSignUp) {
    githubBtnSignUp.addEventListener("click", (event) => {
      event.preventDefault();
      signInTroughGithub();
    });
  }
  if (verifyEmailBtn) {
    verifyEmailBtn.addEventListener("click", (event) => {
      event.preventDefault();
      verifyEmailWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
    });
  }
  if (newChatRoomButtonOnMainContentWrapper) {
    newChatRoomButtonOnMainContentWrapper.addEventListener("click", (event) => {
      event.preventDefault();
      mainContentWrapper.classList.add("signin-portal");

      newChatRoom();
    });
  }
  if (newChatBtnInNewChatRoomWindow) {
    newChatBtnInNewChatRoomWindow.addEventListener("click", (event) => {
      event.preventDefault();
      newChatRoomWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
      createNewChatRoom();
    });
  }
});
