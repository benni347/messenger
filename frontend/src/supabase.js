"use strict";

import {
  RetrieveEnvValues,
  ValidateEmail,
  GenerateUserName,
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

  const { data, error } = await supabase.auth.signIn({
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

  console.info(data);
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
 * Changes the visibility of the sign-in and sign-out buttons based on the 'authenticated' flag in local storage.
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

window.addEventListener("DOMContentLoaded", () => {
  changeButton();
  const signInWindow = document.getElementById("signin-window");
  const signUpWindow = document.getElementById("signup-window");
  const verifyEmailWindow = document.getElementById("verify-window");
  const verifyEmailBtn = document.getElementById("okay-verify-btn");
  const mainContentWrapper = document.getElementById("main-window-wrappper");
  const signOutBtn = document.getElementById("signout-main-wrapper");
  const signInBtn = document.getElementById("signin-button");
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
  const signUpBtn = document.getElementById("signup-btn");
  if (signUpBtn) {
    signUpBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signUp();
      signUpWindow.classList.add("hidden");
      verifyEmailWindow.classList.remove("hidden");
    });
  }
  const closeSignInBtn = document.getElementById("close-signin-button");
  if (closeSignInBtn) {
    closeSignInBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signInWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
      signInWindow.classList.add("hidden");
    });
  }
  const closeSignUpBtn = document.getElementById("close-signup-button");
  if (closeSignUpBtn) {
    closeSignUpBtn.addEventListener("click", (event) => {
      event.preventDefault();
      signUpWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
    });
  }
  const signUpBtnInSignInWindow = document.getElementById("signup");
  if (signUpBtnInSignInWindow) {
    signUpBtnInSignInWindow.addEventListener("click", (event) => {
      event.preventDefault();
      signInWindow.classList.add("hidden");
      signUpWindow.classList.remove("hidden");
    });
  }
  const signInBtnInSignUpWindow = document.getElementById("signin");
  if (signInBtnInSignUpWindow) {
    signInBtnInSignUpWindow.addEventListener("click", (event) => {
      event.preventDefault();
      signUpWindow.classList.add("hidden");
      signInWindow.classList.remove("hidden");
    });
  }
  const signInBtnInMainContentWrapper = document.getElementById(
    "signin-main-wrapper"
  );
  if (signInBtnInMainContentWrapper) {
    signInBtnInMainContentWrapper.addEventListener("click", (event) => {
      event.preventDefault();
      mainContentWrapper.classList.add("signin-portal");
      signInWindow.classList.remove("hidden");
    });
  }
  const githubBtn = document.getElementById("github-button");
  if (githubBtn) {
    githubBtn.addEventListener("click", (event) => {
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
});
