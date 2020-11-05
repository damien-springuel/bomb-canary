<script lang="ts">
import type { Message } from "../messages/messagebus";
import type { StoreValues } from "../store/store";
import { GamePhase } from "../store/store";
import TeamSelection from "./game-phases/TeamSelection.svelte";
import TeamVote from "./game-phases/TeamVote.svelte";
export let storeValues: StoreValues;
export let dispatcher: {dispatch: (message: Message) => void};

</script>

<div class="flex flex-col items-center h-full bg-gray-900 p-6 text-blue-500 space-y-4 text-2xl">
  <h2 class="text-5xl">
    Mission Track
  </h2>
  <div class="flex flex-row justify-around w-full text-gray-900 text-3xl">
    {#each storeValues.missionRequirements as requirement}
      <div class="rounded-full bg-blue-400 h-16 w-16 flex items-center justify-center">
        {requirement.nbPeopleOnMission}
      </div>
    {/each}
  </div>
  {#if storeValues.currentGamePhase === GamePhase.TeamSelection}
    <TeamSelection dispatcher={dispatcher} storeValues={storeValues}/>
  {:else if storeValues.currentGamePhase === GamePhase.TeamVote}
    <TeamVote dispatcher={dispatcher} storeValues={storeValues}/>
  {/if}
</div>