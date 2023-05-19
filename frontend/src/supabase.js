"use strict";

import { RetrieveEnvValues } from "../wailsjs/go/main/App.js";

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

  if (error) {
    console.error(`An error occured during the login: ${error}`);
  }

  console.info(data);
}

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
  if (error !== "") {
    console.error(`An error occured during the logout: ${error}`);
  }
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
  if (error !== "") {
    console.error(`An error occured during the login: ${error}`);
  }
  console.info(data);
}

// Solution for the error:
// https://stackoverflow.com/a/74261887
// -FIX-ME-: Uncaught ReferenceError: createClient is not defined
// Tried fixes:
//  - importing the createClient function from the supabase-js module directly via the cdn in this file
//  - importing it via "import { createClient } from "@supabase/supabase-js";"

RetrieveEnvValues().then((env) => {
  console.log(`Env: ${env}`);
  supabaseKey = env.supaBaseApiKey;
  supabaseUrl = env.supaBaseUrl;
  console.info(`Supabase Key: ${supabaseKey}`);
  console.info(`Supabase URL: ${supabaseUrl}`);
  console.info(`Supabase options: ${JSON.stringify(options)}`);
  supabase = createClient(supabaseUrl, supabaseKey, options);
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
  const email = document.getElementById("email-input").value;
  const password = document.getElementById("password-input").value;
  const user_name = document.getElementById("username").value;

  const { data, error } = await supabase.auth.signUp({
    email, // equals to email: email
    password, // equals to password: password
    options: {
      data: {
        user_name, // equals to user_name: user_name
      },
    },
  });
  if (error !== "") {
    console.error(`An error occured during the creation of the user: ${error}`);
  }

  console.info(data);
}

window.addEventListener("DOMContentLoaded", () => {
  const signInWindow = document.getElementById("signin-window");
  const signUpWindow = document.getElementById("signup-window");
  const mainContentWrapper = document.getElementById("main-window-wrappper");
  const signInBtn = document.getElementById("signin-button");
  if (signInBtn) {
    signInBtn.addEventListener("click", (event) => {
      event.preventDefault();
      console.info("Sign in button clicked");
      signInThroughMail();
    });
  }
  const closeSignInBtn = document.getElementById("close-signin-button");
  if (closeSignInBtn) {
    closeSignInBtn.addEventListener("click", (event) => {
      event.preventDefault();
      console.info("Close sign in button clicked");
      signInWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
      signInWindow.classList.add("hidden");
    });
  }
  const closeSignUpBtn = document.getElementById("close-signup-button");
  if (closeSignUpBtn) {
    closeSignUpBtn.addEventListener("click", (event) => {
      event.preventDefault();
      console.info("Close sign up button clicked");
      signUpWindow.classList.add("hidden");
      mainContentWrapper.classList.remove("signin-portal");
    });
  }
  const signUpBtnInSignInWindow = document.getElementById("signup");
  if (signUpBtnInSignInWindow) {
    signUpBtnInSignInWindow.addEventListener("click", (event) => {
      event.preventDefault();
      console.info("Sign up button in sign in window clicked");
      signInWindow.classList.add("hidden");
      signUpWindow.classList.remove("hidden");
    });
  }
  const signInBtnInSignUpWindow = document.getElementById("signin");
  if (signInBtnInSignUpWindow) {
    signInBtnInSignUpWindow.addEventListener("click", (event) => {
      event.preventDefault();
      console.info("Sign in button in sign up window clicked");
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
      console.info("Sign in button in main content wrapper clicked");
      mainContentWrapper.classList.add("signin-portal");
      signInWindow.classList.remove("hidden");
    });
  }
  const githubBtn = document.getElementById("github-button");
  if (githubBtn) {
    githubBtn.addEventListener("click", (event) => {
      event.preventDefault();
      console.info("Github button clicked");
      signInTroughGithub();
    });
  }
});
