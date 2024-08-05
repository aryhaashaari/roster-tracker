<template>
  <div class="max-w-screen-sm mx-auto px-4 py-10">
    <!-- Error Handling -->
     <div v-if="errorMsg" class="mb-10 p-4 rounded-md bg-grey">
      <p class="text-red-500">{{ errorMsg }}</p>
     </div>

     <!-- Login -->
      <form @submit.prevent="login" class="p-8 flex flex-col bg-grey rounded-md shadow-lg">
        <h1 class="text-3xl text-black mb-4">Login</h1>

        <div class="flex flex-col mb-2">
          <label for="email" class="mb-1 text-sm text-black">Email</label>
          <input type="text" required class="p-2 text-gray-500 focus:outline-none" id="email" v-model="email">
        </div>

        <div class="flex flex-col mb-2">
          <label for="password" class="mb-1 text-sm text-black">Password</label>
          <input type="password" required class="p-2 text-gray-500 focus:outline-none" id="password" v-model="password">
        </div>

        <button type="submit" class="mt-6 py-2 px-6 rounded-sm self-start text-sm text-white bg-black duration-200 
        border-solid border-2 border-transparent hover:border-black hover:bg-white hover:text-black">Login</button>

        <router-link class="text-sm mt-6 text-center" :to="{name: 'Register'}">
          Don't have an account? <span class="text-white">Register</span>
        </router-link>
      </form>
  </div>
</template>

<script>
import { ref } from "vue";
import { supabase } from "../supabase/init";
import { useRouter } from "vue-router";
import axios from 'axios';
import moment from 'moment';
import CryptoJS from 'crypto-js';

export default {
  name: "login",
  setup() {
    // Create data / vars
    const router = useRouter();
    const email = ref(null);
    const password = ref(null);
    const errorMsg = ref(null);
    var token = ref(null);

    const generateTimestamp = () => {
      return moment().format("YYYY-MM-DDTHH:mm:ssZ");
    };

    const generateSignature = (timestamp, method, json_body_raw, x_api_key_id, x_api_key_secret) => {
      if (json_body_raw && Object.keys(JSON.parse(json_body_raw)).length !== 0) {
        json_body_raw = json_body_raw.replace(/\n/g, "").replace(/\\/g, "").replace(/\s/g, "");
      } else {
        json_body_raw = "";
      }

      const body_md5 = CryptoJS.MD5(json_body_raw).toString(CryptoJS.enc.Base64);
      const hmac_signature = `${timestamp}:${x_api_key_id}:${method}:${body_md5}`;
      const hmac = CryptoJS.HmacSHA256(hmac_signature, x_api_key_secret);
      const hmac_base64 = hmac.toString(CryptoJS.enc.Base64);
      let signature = `#${x_api_key_id}:#${hmac_base64}`;
      return btoa(signature);
    };


    // Login function
    const login = async() => {
      try {
        const {error} = await supabase.auth.signIn({
          email: email.value,
          password: password.value,
        });
        if (error) throw error;

        // Generate timestamp and signature
        const timestamp = generateTimestamp();
        const method = 'GET'; // Assuming GET request
        const json_body_raw = ''; // Assuming empty body for GET request
        const x_api_key_id = '859858c7-cef1-48fe-ba81-2d66f6ae8998'; // Replace with actual API key ID
        const x_api_key_secret = 'aryha3171021509031101'; // Replace with actual API key secret
        const signature = generateSignature(timestamp, method, json_body_raw, x_api_key_id, x_api_key_secret);

        // Make GET request with axios
        const response = await axios.get('http://localhost:8080/login', {
          headers: {
            'X-Api-Key-Id': x_api_key_id,
            'Timestamp': timestamp,
            'Signature': signature
          }
        });

        token =  response.data.data.token
        console.log('API response:', response.data);
        router.push({name: "Home"});
      } catch (error) {
        errorMsg.value = `Error: ${error.message}`;
        setTimeout(() => {
        errorMsg.value = null;
      }, 5000);
      }

  
    }

    return {token, email, password, errorMsg, login};
  },
};
</script>
