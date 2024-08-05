<template>
  <div class="max-w-screen-sm mx-auto px-4 py-10">
    <!-- App Msg -->
     <div v-if="statusMsg || errorMsg" class="mb-10 p-4 rounded-md shadow-md bg-grey">
      <p class="text-green-500">
        {{statusMsg}}
      </p>
      <p class="text-red-500">
        {{errorMsg}}
      </p>
     </div>

    <div v-if="dataLoaded">
      <!-- General Player Info -->
       <div class="flex flex-col items-center p-8 rounded-md shadow-md bg-grey relative">
        <div v-if="user" class="flex absolute left-3 top-3 gap-x-2">
          <div @click="editMode" class="h-7 w-7 rounded-full flex justify-center items-center cursor-pointer bg-white shadow-lg">
            <img class="h-3.5 w-auto" src="@/assets/images/pencil.png" alt=""/>
          </div>
          <div @click="deletePlayer" class="h-7 w-7 rounded-full flex justify-center items-center cursor-pointer bg-white shadow-lg">
            <img class="h-4 w-auto" src="@/assets/images/trash.png" alt=""/>
          </div>
        </div>

        <img class="h-24 w-auto" src="@/assets/images/person.png" alt=""/>

        <div class="mt-3">
          <!-- <input v-if="edit" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="data.position"/> -->
          <select v-if="edit" type="text" class="p-2 text-center text-black focus:outline-none" id="position" v-model="data.position">
            <option>Select Position</option>
            <option value="G">Guard</option>
            <option value="F">Forward</option>
            <option value="C">Center</option>
          </select>
          <span v-else class="mt-6 py-1.5 px-5 text-xs text-white bg-black rounded-lg shadow-md">
          {{data.position}}
          </span>
        </div>

        

        <div class="w-full mt-6">
          <input v-if="edit" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="data.player_name"/>
          <h1 v-else class="text-black text-2xl text-center">
            {{data.player_name}}
          </h1>
        </div>
       </div>
      
      <!-- Exercises -->
       <div class="mt-10 p-8 rounded-md flex flex-col item-center bg-grey shadow-md">
        <!-- Physical Information -->
        
         <div class="flex flex-col gap-y-2 w-full">
          <p class="text-l">Physical Information</p>
          <div class="flex flex-col gap-x-6 gap-y-2 relative sm:flex-row">
            <div class="flex flex-1 flex-col">
              <label for="exercise-name" class="mb-1 text-xs text-black">
                Height (cm)
              </label>
              <input v-if="edit" id="exercise-name" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="data.physique.height"/>
              <p v-else>
                {{data.physique.height}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="sets" class="mb-1 text-xs text-black">
                Weight (kg)
              </label>
              <input v-if="edit" id="sets" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="data.physique.weight"/>
              <p v-else>
                {{data.physique.weight}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="reps" class="mb-1 text-xs text-black">
                Age
              </label>
              <input v-if="edit" id="reps" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="data.physique.age"/>
              <p v-else>
                {{data.physique.age}}
              </p>
            </div>
            
            <div class="flex flex-1 flex-col">
              <label for="weight" class="mb-1 text-xs text-black">
                Wingspan
              </label>
              <input v-if="edit" id="weight" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="data.physique.wingspan"/>
              <p v-else>
                {{data.physique.wingspan}}
              </p>
            </div>
            
          </div>
         </div>

        <!-- Game Stats -->
        <div v-if="data.stats !== [1]" class="flex flex-col gap-y-2 w-full mt-7">
          <p class="text-l">Game Stats</p>
          <div v-if="data.stats.length === 0">
            <p class="text-xs">No data available for this player</p>
          </div>
          <div v-for="(item, index) in data.stats" :key="index" class="flex flex-col gap-x-6 gap-y-2 relative sm:flex-row">
            <div class="flex flex-1 flex-col">
              <label for="points" class="mb-1 text-xs text-black">
                Points
              </label>
              <input v-if="edit" id="points" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.points"/>

              <p v-else>
                {{item.points}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="assists" class="mb-1 text-xs text-black">
                Assists
              </label>
              <input v-if="edit" id="assists" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.assists"/>
              <p v-else>
                {{item.assists}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="rebounds" class="mb-1 text-xs text-black">
                Rebounds
              </label>
              <input v-if="edit" id="rebounds" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.rebounds"/>
              <p v-else>
                {{item.rebounds}}
              </p>
            </div>
            
            <div class="flex flex-1 flex-col">
              <label for="field-goal" class="mb-1 text-xs text-black">
                FG%
              </label>
              <input v-if="edit" id="field-goal" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.fieldgoalpct"/>
              <p v-else>
                {{item.fieldgoalpct}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="three-point" class="mb-1 text-xs text-black">
                3P%
              </label>
              <input v-if="edit" id="three-point" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.threepointpct"/>
              <p v-else>
                {{item.threepointpct}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="steals" class="mb-1 text-xs text-black">
                Steals
              </label>
              <input v-if="edit" id="steals" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.steals"/>
              <p v-else>
                {{item.steals}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="blocks" class="mb-1 text-xs text-black">
                Blocks
              </label>
              <input v-if="edit" id="blocks" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.blocks"/>
              <p v-else>
                {{item.blocks}}
              </p>
            </div>

            <div class="flex flex-1 flex-col">
              <label for="turnovers" class="mb-1 text-xs text-black">
                Turnovers
              </label>
              <input v-if="edit" id="turnovers" type="text" class="p-2 w-full text-grey-500 focus:outline-none" v-model="item.turnovers"/>
              <p v-else>
                {{item.turnovers}}
              </p>
            </div>

            <img v-if="edit" @click="deleteStats(item.stats_id)" src="@/assets/images/trash.png" class="absolute h-4 w-auto -left-5 cursor-pointer" alt=""/>
            
          </div>

          <button @click="addStats" v-if="edit" type="button" class="mt-6 py-2 px-6 rounded-sm self-start text-sm text-black bg-white duration-200 
          border-solid border-2 border-transparent hover:border-black hover:bg-white hover:text-black">
          Add Stats
          </button>
         </div>

         <!-- Player Stats Summary -->
         <!-- <div v-if="!edit && data.stats.length !== 0" class="flex flex-col gap-y-2 w-full mt-7">
          <p class="text-l">Stats Summary</p>
          <div class="flex flex-col gap-x-6 gap-y-2 relative sm:flex-row">
            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                Total Games
              </label>
                {{data.stats.length}}
            </div>
          </div>
          <div class="flex flex-col gap-x-6 gap-y-2 relative sm:flex-row">
            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                PPG
              </label>
                {{data.avgStats.points}}
            </div>

            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                APG
              </label>
                {{data.avgStats.assists}}
            </div>

            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                RPG
              </label>
                {{data.avgStats.rebounds}}
            </div>

            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                FG%
              </label>
                {{data.avgStats.fieldGoalPct}}
            </div>

            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                3P%
              </label>
                {{data.avgStats.threePointPct}}
            </div>

            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                SPG
              </label>
                {{data.avgStats.steals}}
            </div>

            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                BPG
              </label>
                {{data.avgStats.blocks}}
            </div>

            <div class="flex flex-1 flex-col">
              <label class="mb-1 text-xs text-black">
                TPG
              </label>
                {{data.avgStats.turnovers}}
            </div>
          </div>
         </div> -->

       </div>

      <!-- Update -->
      <button 
      @click="update"
      v-if="edit"
      type="button" 
      class="mt-6 py-2 px-6 rounded-sm self-start text-sm text-white bg-black duration-200 
      border-solid border-2 border-transparent hover:border-black hover:bg-white hover:text-black">
          Update Player
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import store from '../store';
import { uid } from 'uid';
import axios from 'axios';

export default {
  name: "view-workout",
  setup() {
    // Create data / vars
    const data = ref(null);
    const dataLoaded = ref(null);
    const errorMsg = ref(null);
    const statusMsg = ref(null);
    const route = useRoute();
    const router = useRouter();
    const user = computed(() => store.state.user);
    // Get current Id of route
    const currentId = route.params.playerId;

    // Get workout data
    const getData = async () => {
      try {
            const apiUrl = `http://localhost:8080/getPlayer/${currentId}`;
            const response = await axios.get(apiUrl);
            data.value = response.data.data;
            dataLoaded.value = true;
        // statsAvg();
      } catch (error) {
        errorMsg.value = error.message;
        setTimeout(() => {
          errorMsg.value = false;
        }, 5000);
      }
    }

    getData();
    
    // Get Player Stats Summary
    // const statsAvg = async () => {
    //   try {
    //     data.value.avgStats.points = '0';
    //     data.value.avgStats.assists = '0';
    //     data.value.avgStats.rebounds = '0';
    //     data.value.avgStats.fieldGoalPct = '0';
    //     data.value.avgStats.threePointPct = '0';
    //     data.value.avgStats.steals = '0';
    //     data.value.avgStats.blocks = '0';
    //     data.value.avgStats.turnovers = '0';

    //     for (let i = 0; i < data.value.stats.length; i++){
    //       data.value.avgStats.points = String(parseInt(data.value.avgStats.points) + parseInt(data.value.stats[i].points));
    //       data.value.avgStats.assists = String(parseInt(data.value.avgStats.assists) + parseInt(data.value.stats[i].assists));
    //       data.value.avgStats.rebounds = String(parseInt(data.value.avgStats.rebounds) + parseInt(data.value.stats[i].rebounds));
    //       data.value.avgStats.fieldGoalPct = String(parseInt(data.value.avgStats.fieldGoalPct) + parseInt(data.value.stats[i].fieldGoalPct));
    //       data.value.avgStats.threePointPct = String(parseInt(data.value.avgStats.threePointPct) + parseInt(data.value.stats[i].threePointPct));
    //       data.value.avgStats.steals = String(parseInt(data.value.avgStats.steals) + parseInt(data.value.stats[i].steals));
    //       data.value.avgStats.blocks = String(parseInt(data.value.avgStats.blocks) + parseInt(data.value.stats[i].blocks));
    //       data.value.avgStats.turnovers = String(parseInt(data.value.avgStats.turnovers) + parseInt(data.value.stats[i].turnovers));
    //     }

    //     data.value.avgStats.points = (data.value.avgStats.points/data.value.stats.length).toFixed(1);
    //     data.value.avgStats.assists = (data.value.avgStats.assists/data.value.stats.length).toFixed(1);
    //     data.value.avgStats.rebounds = (data.value.avgStats.rebounds/data.value.stats.length).toFixed(1);
    //     data.value.avgStats.fieldGoalPct = (data.value.avgStats.fieldGoalPct/data.value.stats.length).toFixed(1);
    //     data.value.avgStats.threePointPct = (data.value.avgStats.threePointPct/data.value.stats.length).toFixed(1);
    //     data.value.avgStats.steals = (data.value.avgStats.steals/data.value.stats.length).toFixed(1);
    //     data.value.avgStats.blocks = (data.value.avgStats.blocks/data.value.stats.length).toFixed(1);
    //     data.value.avgStats.turnovers = (data.value.avgStats.turnovers/data.value.stats.length).toFixed(1);
        
    //     const {error} = await supabase
    //     .from('players')
    //     .update({avgStats: data.value.avgStats})
    //     .eq('id', currentId);

    //     if (error) throw error;

    //   } catch (error) {
    //     errorMsg.value = `Error: ${error.message}`;
    //     setTimeout(() => {
    //       errorMsg.value = false;
    //     }, 5000);
    //   }
    // }

    // Delete workout
    const deletePlayer = async () => {
      try {
        const apiUrl = `http://localhost:8080/deletePlayer/${currentId}`
        const {error} = await axios.delete(apiUrl);
        if (error) throw error;
        router.push({name: "Home"});
      } catch (error) {
        errorMsg.value = `Error: ${error.message}`;
        setTimeout(() => {
          errorMsg.value = false;
        }, 5000);
      }
    }

    // Edit mode
    const edit = ref(null);

    const editMode = () => {
      edit.value = !edit.value;
    };

    // Add exercise
    const addStats = () => {
      if(data.value.stats.position !== null) {
        data.value.stats.push({
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
      errorMsg.value = "Error: Cannot add stats";
      setTimeout(() => {
        errorMsg.value = false;
      }, 5000);

    };

    // Delete stats
    const deleteStats = (stats_id) => {
      if(data.value.stats.length > 0){
        data.value.stats = data.value.stats.filter(
          stat => stat.stats_id !== stats_id);
          console.log(data.value.stats);
        return;
      }
      errorMsg.value = "Error: Cannot remove";
      setTimeout(() => {
        errorMsg.value = false;
      }, 5000);
    };

    // Update Workout
    const update = async () => {
      try {
        const apiUrl = `http://localhost:8080/updatePlayer/${currentId}`
        const {error} = await axios.put(apiUrl, {
          id: currentId,
          player_name: data.value.player_name,
          position: data.value.position,
          physique: data.value.physique,
          stats: data.value.stats
        });
        if (error) throw error;
        edit.value = false;
        statusMsg.value = "Success: Player Updated!";
        setTimeout(() => {
          statusMsg.value = false;
        }, 5000);
      } catch (error) {
        errorMsg.value = `Error: ${error.message}`;
        setTimeout(() => {
          errorMsg.value = false;
        }, 5000);
      }
    }

    return {statusMsg, errorMsg, data, dataLoaded, edit, editMode, user, deletePlayer, addStats, deleteStats, update};
  },
};
</script>
