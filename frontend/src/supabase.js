"use strict";

import { RetrieveEnvValues } from "../wailsjs/go/main/App.js";

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
const supabaseUrl = RetrieveEnvValues().then((env) => {
  return env.supaBaseUrl;
});
const supabaseKey = RetrieveEnvValues().then((env) => {
  return env.supaBaseApiKey;
});

// FIXME: Uncaught ReferenceError: createClient is not defined
// Tried fixes:
//  - importing the createClient function from the supabase-js module directly via the cdn in this file
//  - importing it via "import { createClient } from "@supabase/supabase-js";"
const supabase = createClient(supabaseUrl, supabaseKey, options);

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
