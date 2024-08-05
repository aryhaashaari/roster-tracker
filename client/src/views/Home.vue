<template>
    <div v-if="dataLoaded" class="container mt-10 px-4">
      <!-- No Data -->
       <div v-if="data.length === 0" class="w-full flex flex-col items-center">
        <h1 class="text-2xl">Looks empty here...</h1>
        <router-link :to="{name: 'Create'}" class="mt-6 py-2 px-6 rounded-sm text-sm text-black bg-grey duration-200 
        border-solid border-2 border-transparent hover:border-black hover:bg-white hover:text-black">
        Add Player
        </router-link>
       </div>
  
      <!-- Data -->
      <div v-else class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-6">
        <router-link class="flex flex-col items-center bg-grey p-8 shadow-md cursor-pointer"
        :to="{name: 'View-Player', params: {playerId: player.id}}"
        v-for="(player, index) in data"
        :key="index"
        >
          <!-- Player Card -->
          <img src="@/assets/images/person.png" class="h-24 w-auto" alt=""/>
  
          <p class="mt-6 py-1 px-3 text-xs text-white bg-black shadow-md rounded-lg">
            {{player.position}}
          </p>
  
          <h1 class="mt-8 mb-2 text-center text-xl text-black">
            {{player.player_name}}
          </h1>
        </router-link>
      </div>
    </div>
  </template>
  
  <script>
  import { ref } from "vue";
  import axios from "axios";
  
  export default {
    name: "Home",
    components: {},
    setup() {
      // Create data / vars
      const data = ref([]);
      const dataLoaded = ref(null);
  
      // Get data
      const getData = async () => {
        try {
            const response = await axios.get("http://localhost:8080/getPlayers");
            data.value = response.data.data;
            dataLoaded.value = true;
        } catch (error) {
            this.errorMessage = error.message;
            console.warn("There was an error!", error);
        }
      }
  
      // Run data function
      getData();
  
      return {data, dataLoaded};
    },
  };
  </script>