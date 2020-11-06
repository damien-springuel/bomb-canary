<script lang="ts">
import { ApproveTeam, RejectTeam } from "../../messages/commands";

  import type { Message } from "../../messages/messagebus";
  import type { StoreValues } from "../../store/store";
  export let storeValues: StoreValues;
  export let dispatcher: {dispatch(message: Message): void};

  function approve(){
    dispatcher.dispatch(new ApproveTeam());
  }

  function reject() {
    dispatcher.dispatch(new RejectTeam());
  }

  </script>
  
  <div class="flex flex-col items-center h-full w-full">
    <div class="text-3xl">
      Team
    </div>
    <div class="flex flex-row flex-wrap justify-center gap-x-2 w-full text-blue-500">
      {#each Array.from(storeValues.currentTeam.values()) as member, i}
        {#if i !== 0}
          <div>|</div>
        {/if}
        <div>{member}</div>
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