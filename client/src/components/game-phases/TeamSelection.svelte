<script lang="ts">
  import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember } from "../../messages/commands";
  import type { Message } from "../../messages/messagebus";
  import type { StoreValues } from "../../store/store";
  export let storeValues: StoreValues;
  export let dispatcher: {dispatch: (message: Message) => void};
  
  function togglePlayerSelection(member: string): void {
    if (storeValues.isPlayerTheLeader) {
      if (storeValues.isPlayerInTeam(member)) {
        dispatcher.dispatch(new LeaderDeselectsMember(member));
      } else {
        dispatcher.dispatch(new LeaderSelectsMember(member));
      }
    }
  }
  
  function confirmTeam(): void {
    dispatcher.dispatch(new LeaderConfirmsTeam());
  }
  </script>
  
  <div class="flex flex-col items-center h-full w-full">
    <div class="text-xl">
      <span class="font-bold">{storeValues.leader}</span> is choosing current team.
    </div>
    <div class="flex-grow grid grid-cols-2 w-full content-start gap-2">
      {#each storeValues.players as player}
        <button 
          class="bc-button bc-button-blue" 
          class:bc-button-green={storeValues.isPlayerInTeam(player)} 
          on:click={() => togglePlayerSelection(player)}
          disabled={!storeValues.isPlayerSelectableForTeam(player)}
          class:bc-button-gray={!storeValues.isPlayerSelectableForTeam(player)}
        >
          {player}
        </button>
      {/each}
    </div>
    {#if storeValues.isPlayerTheLeader}
      <button 
        class="bc-button bc-button-blue" 
        on:click={confirmTeam}
        disabled={!storeValues.canConfirmTeam}
        class:bc-button-gray={!storeValues.canConfirmTeam}
      >
        I'm done
      </button>
    {/if}
  </div>