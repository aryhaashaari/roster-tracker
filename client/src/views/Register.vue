<template>
  <div class="max-w-screen-sm mx-auto px-4 py-10">
    <!-- Error Handling -->
     <div v-if="errorMsg" class="mb-10 p-4 rounded-md bg-light-grey">
      <p class="text-red-500">{{ errorMsg }}</p>
     </div>

     <!-- Register -->
      <form @submit.prevent="register" class="p-8 flex flex-col bg-grey rounded-md shadow-lg">
        <h1 class="text-3xl text-black mb-4">Register</h1>

        <div class="flex flex-col mb-2">
          <label for="email" class="mb-1 text-sm text-black">Email</label>
          <input type="text" required class="p-2 text-black focus:outline-none" id="email" v-model="email">
        </div>

        <div class="flex flex-col mb-2">
          <label for="password" class="mb-1 text-sm text-black">Password</label>
          <input type="password" required class="p-2 text-black focus:outline-none" id="password" v-model="password">
        </div>

        <div class="flex flex-col mb-2">
          <label for="confirmPassword" class="mb-1 text-sm text-black">Confirm Password</label>
          <input type="password" required class="p-2 text-black focus:outline-none" id="confirmPassword" v-model="confirmPassword">
        </div>

        <div class="flex flex-col mb-2">
          <label for="role" class="mb-1 text-sm text-black">Role</label>
          <select type="text" required @change="positionChange" class="p-2 text-black focus:outline-none" id="role" v-model="code">
            <option>Select Role</option>
            <option value="Coach">Coach</option>
            <option value="Player">Player</option>
          </select>
        </div>

        <button type="submit" class="mt-6 py-2 px-6 rounded-sm self-start text-sm text-white bg-black duration-200 
        border-solid border-2 border-transparent hover:border-black hover:bg-white hover:text-black">Register</button>

        <router-link class="text-sm mt-6 text-center" :to="{name: 'Login'}">
          Already have an account? <span class="text-white">Login</span>
        </router-link>
      </form>
  </div>
</template>

<script>
import {ref} from "vue";
import {supabase} from "../supabase/init.js";
import {useRouter} from "vue-router";
import axios from "axios";

export default {
  name: "register",
  setup() {
    // Create data / vars
    const router = useRouter();
    const email = ref(null);
    const password = ref(null);
    const confirmPassword = ref(null);
    const role = ref(null);
    const errorMsg = ref(null);
    const code = ref(null);

    // Register function
    const register = async() => {
      if(password.value === confirmPassword.value) {
        try {
          const {error} = await axios.post('http://localhost:8080/role/register', {
          email: email.value,
          password: password.value,
          code: code.value,
          });
          const {error2} = await supabase.auth.signUp({
            email: email.value,
            password: password.value,
          });
          if (error) throw error;
          if (error2) throw error2;
          router.push({ name: "Login"}); 
        } catch (error) {
          errorMsg.value = error.message;
          setTimeout(() => {
            errorMsg.value = null;
          }, 5000);
          }
        return; 
      }
      errorMsg.value = "Error: Passwords do not match";
      setTimeout(() => {
        errorMsg.value = null;
      }, 5000);
    };

    return {email, password, confirmPassword, role, errorMsg, register};
  },
};
</script>
