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
