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
    <div class="flex flex-row w-full text-blue-500 text-center ">
      {#each Array.from(storeValues.currentTeam.values()) as member, i}
        {#if i !== 0}
          <div>|</div>
        {/if}
        <div class="flex-grow w-full">{member}</div>
      {/each}
    </div>
    <div class="grid grid-cols-2 justify-center w-full gap-x-2 mt-4">
      <button 
        class="bc-button bc-button-green" 
        on:click={approve}
        class:text-gray-700={storeValues.hasGivenPlayerVoted(storeValues.player)}
        class:bg-gray-500={storeValues.playerVote !== true}
        class:hover:bg-gray-500={storeValues.playerVote !== true}
        disabled={storeValues.hasGivenPlayerVoted(storeValues.player)}
      >
        Approve
      </button>
      <button 
        class="bc-button bc-button-red" 
        on:click={reject}
        class:text-gray-700={storeValues.hasGivenPlayerVoted(storeValues.player)}
        class:bg-gray-500={storeValues.playerVote !== false}
        class:hover:bg-gray-500={storeValues.playerVote !== false}
        disabled={storeValues.hasGivenPlayerVoted(storeValues.player)}
      >
        Reject
      </button>
    </div>
    
    <div class="text-3xl mt-8">
      Votes
    </div>
    <div class="flex-grow grid grid-cols-3 w-full content-start gap-2 text-lg text-center">
      {#each storeValues.players as player}
        <div 
          class="border border-blue-500 rounded-lg p-1" 
          class:bg-blue-500={storeValues.hasGivenPlayerVoted(player)}
          class:text-gray-900={storeValues.hasGivenPlayerVoted(player)}
        >
          {player}
        </div>
      {/each}
    </div>
  </div>