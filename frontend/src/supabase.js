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
  global: {
    headers: { "x-my-custom-header": "my-app-name" },
  },
};

let supabaseKey;
let supabaseUrl;
let supabase;
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

async function signUp(mail, password) {
  const { data, error } = await supabase.auth.signUp({
    email: mail,
    password: password,
  });
  if (error != "") {
    console.error(`An error occured during the creation of the user: ${error}`);
  }

  console.info(data);
}
