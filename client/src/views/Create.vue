<template>
  <div class="max-w-screen-md mx-auto px-4 py-10">
    <!-- Status Message -->
    <div v-if="statusMsg || errorMsg" class="mb-10 p-4 bg-grey rounded-md shadow-lg">
      <p class="text-green-500">
        {{statusMsg}}
      </p>
      <p class="text-red-500">
        {{errorMsg}}
      </p>
    </div>

    <!-- Create -->
    <div class="p-8 flex items-start bg-grey rounded-md shadow-lg">
      <!-- Form -->
       <form @submit.prevent="createPlayer" class="flex flex-col gap-y-5 w-full">
        <h1 class="text-2xl text-black">Add Player</h1>

        <!-- Player Name -->
        <div class="flex flex-col">
          <label for="player-name" class="mb-1 text-sm text-black">Player Name</label>
          <input type="text" required class="p-2 text-gray-500 focus:outline-none" id="player-name" v-model="player_name"/>
        </div>

        <!-- Player Position -->
        <div class="flex flex-col">
          <label for="position" class="mb-1 text-sm text-black">Player Position</label>
          <select type="text" required @change="positionChange" class="p-2 text-gray-500 focus:outline-none" id="position" v-model="position">
            <option>Select Position</option>
            <option value="G">Guard</option>
            <option value="F">Forward</option>
            <option value="C">Center</option>
          </select>
        </div>
        
        <!-- Physical Information Input -->
        <div v-if="position === 'G' || position === 'F' || position === 'C'" class="flex flex-col gap-y-4">
          <div class="flex flex-col gap-x-6 gap-y-2 relative md:flex-row" 
          >
            <div class="flex flex-col flex-1">
              <label for="height" class="mb-1 text-sm text-black">Height (cm)</label>
              <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="physique.height"/>
            </div>
            <div class="flex flex-col flex-1">
              <label for="weight" class="mb-1 text-sm text-black">Weight (kg)</label>
              <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="physique.weight"/>
            </div>
            <div class="flex flex-col flex-1">
              <label for="age" class="mb-1 text-sm text-black">Age</label>
              <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="physique.age"/>
            </div>
            <div class="flex flex-col flex-1">
              <label for="wingspan" class="mb-1 text-sm text-black">Wingspan (cm)</label>
              <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="physique.wingspan"/>
            </div>
          </div>
        </div>

        <!-- Game Stats Inputs -->
        <div v-if="position === 'G' || position === 'F' || position === 'C'" class="flex flex-col gap-y-4">
          <div 
          class="flex flex-col gap-x-6 gap-y-2 relative md:flex-row" 
          v-for="(item, index) in stats"
          :key="index"
          >
          <div class="flex flex-col gap-y-4">
            <div class="flex flex-row gap-x-6">
              <div class="flex flex-col flex-1">
                <label for="points" class="mb-1 text-sm text-black">Points</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.points"/>
              </div>
              <div class="flex flex-col flex-1">
                <label for="assists" class="mb-1 text-sm text-black">Assists</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.assists"/>
              </div>
              <div class="flex flex-col flex-1">
                <label for="rebounds" class="mb-1 text-sm text-black">Rebounds</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.rebounds"/>
              </div>
              <div class="flex flex-col flex-1">
                <label for="fieldGoalPct" class="mb-1 text-sm text-black">FG%</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.fieldGoalPct"/>
              </div>
            </div>
          

            <div class="flex flex-row gap-x-6">
              <div class="flex flex-col flex-1">
                <label for="exercise-name" class="mb-1 text-sm text-black">3P%</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.threePointPct"/>
              </div>
              <div class="flex flex-col flex-1">
                <label for="sets" class="mb-1 text-sm text-black">Steals</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.steals"/>
              </div>
              <div class="flex flex-col flex-1">
                <label for="reps" class="mb-1 text-sm text-black">Blocks</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.blocks"/>
              </div>
              <div class="flex flex-col flex-1">
                <label for="weight" class="mb-1 text-sm text-black">Turnovers</label>
                <input type="text" required class="p-2 w-full text-gray-500 focus:outline-none" v-model="item.turnovers"/>
              </div>
            </div>
          </div>
            <img @click="deleteStats(item.stats_id)" src="@/assets/images/trash.png" class="h-5 w-auto absolute -left-6 cursor-pointer" alt="" />
          </div>
          <button @click="addStats" type="button" class="mt-6 py-2 px-6 rounded-sm self-start text-sm text-black bg-white duration-200 
            border-solid border-2 border-transparent hover:border-black hover:bg-white hover:text-black">Add Game Stats</button>
        </div>

          <button type="submit" class="mt-6 py-2 px-6 rounded-sm self-start text-sm text-white bg-black duration-200 
            border-solid border-2 border-transparent hover:border-black hover:bg-white hover:text-black">
            Create Player
          </button>

       </form>
    </div>
  </div>
</template>

<script>
import { ref } from "vue";
import { uid } from "uid";
import axios from "axios";

export default {
  name: "create",
  setup() {
    // Create data
    const player_name = ref('');
    const position = ref('select-position');
    const physique = ref({
      height: '',
      weight: '',
      age: '',
      wingspan: ''
    });
    const stats = ref([1]);
    const statusMsg = ref(null);
    const errorMsg = ref(null);

    // const avgStats = ref({
    //   points: '0',
    //   assists: '0',
    //   rebounds: '0',
    //   fieldGoalPct: '0',
    //   threePointPct: '0',
    //   steals: '0',
    //   blocks: '0',
    //   turnovers: '0',
    // })

    // Add game stats
    const addStats = () => {
      if(position.value !== null) {
        stats.value.push({
          stats_id: uid(),
          points: '',
          assists: '',
          rebounds: '',
          fieldGoalPct: '',
          threePointPct: '',
          steals: '',
          blocks: '',
          turnovers: '',
        });
        return
      }
      errorMsg.value = "Error: Unable to add player, position is empty";
      setTimeout(() => {
        errorMsg.value = false;
      }, 5000);
    };

    // Delete stats
    const deleteStats = (id) => {
      if(stats.value.length > 0){
        stats.value = stats.value.filter(player => player.stats_id !== id);
        return;
      }
      errorMsg.value = "Error: Cannot remove";
      setTimeout(() => {
        errorMsg.value = false;
      }, 5000);
    };

    // Listens for changing of player position input
    const positionChange = () => {
      stats.value = [];
      addStats();
    };

    // Create player
    const createPlayer = async () => {
      try {
        const {error} = await axios.post('http://localhost:8080/addPlayer', {
          player_name: player_name.value,
          position: position.value,
          physique: physique.value,
          stats: stats.value
        });
        console.log(stats.value)
        if (error) throw error;
        statusMsg.value = 'Success: Player Added!';
        player_name.value = null;
        position.value = "select-position";
        physique.value = {};
        stats.value = [];
        setTimeout (() => {
          statusMsg.value = false;
        }, 5000);
      } catch (error) {
        errorMsg.value = `Error sending stats: ${error.message}`;
          setTimeout(() => {
          errorMsg.value = false;
        }, 5000);
      }
    };

    return {player_name, position, physique, stats, statusMsg, errorMsg, addStats, deleteStats, positionChange, createPlayer};
  },
};
</script>