<script lang="ts">
  import type { Message } from "../../messages/messagebus";
  import type { StoreValues } from "../../store/store";
  export let storeValues: StoreValues;
  export let dispatcher: {dispatch: (message: Message) => void};

  function approve(){
    console.log(`approve team`);
  }

  function reject() {
    console.log(`reject team`);
  }

  </script>
  
  <div class="flex flex-col items-center h-full w-full">
    <div class="text-3xl">
      Team
    </div>
    <div class="flex flex-row flex-wrap justify-around gap-x-2 w-full">
      {#each Array.from(storeValues.currentTeam.values()) as member}
        <div class="text-blue-500">
          {member}
        </div>
      {/each}
    </div>
    <div class="flex flex-row justify-center w-full gap-x-2 mt-4">
      <button class="bc-button bc-button-green" on:click={approve}>
        Approve
      </button>
      <button class="bc-button bc-button-red" on:click={reject}>
        Reject
      </button>
    </div>
  </div>