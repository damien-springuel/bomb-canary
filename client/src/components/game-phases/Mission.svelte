<script lang="ts">
  import type { Message } from "../../messages/messagebus";
  import type { StoreValues } from "../../store/store";
  export let storeValues: StoreValues;
  export let dispatcher: {dispatch(message: Message): void};

  function succeed(){
    console.log(`Succeed mission`);
  }

  function fail() {
    console.log(`Fail mission`);
  }

  </script>
  
  <div class="flex flex-col items-center h-full w-full">
    <div class="text-3xl mt-8">
      Conducting Mission
    </div>
    <div class="grid grid-cols-2 w-full content-start gap-2 text-lg text-center">
      {#each Array.from(storeValues.currentTeam.values()) as member}
        <div class="border border-blue-500 rounded-lg p-1" >
          {member}
        </div>
      {/each}
    </div>
    <div class="flex-grow content-center flex flex-col justify-center">
      <div class="text-4xl text-center">You are part of the mission!</div>
      <div class="grid grid-cols-2 justify-center w-full gap-x-2 mt-4">
        <button 
          class="bc-button bc-button-green" 
          on:click={succeed}
        >
          Succeed
        </button>
        <button 
          class="bc-button bc-button-red" 
          on:click={fail}
        >
          Fail
        </button>
      </div>
    </div>
  </div>