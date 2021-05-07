<script lang="ts">
import { ViewIdentity } from "../messages/commands";
import type { StoreValues } from "../store/store";
import { GamePhase, Dialog as DialogValues } from "../store/store";
import Mission from "./game-phases/Mission.svelte";
import TeamSelection from "./game-phases/TeamSelection.svelte";
import TeamVote from "./game-phases/TeamVote.svelte";
import Identity from './Identity.svelte';
import Dialog from './Dialog.svelte';
import type { Dispatcher } from "../messages/dispatcher";
export let storeValues: StoreValues;
export let dispatcher: Dispatcher;

function viewIdentity() {
  dispatcher.dispatch(new ViewIdentity());
}
</script>

<div class="flex flex-col items-center h-full bg-gray-900 p-6 text-blue-500 space-y-4 text-2xl">
  <button class="bc-button bc-button-blue" on:click={viewIdentity}>Identity</button>
  <div class="text-5xl">
    Mission Track
  </div>
  <div class="flex flex-row justify-around w-full text-3xl">
    {#each storeValues.missionRequirements as requirement, i}
      <div 
        class="rounded-full border border-blue-400 h-16 w-16 flex items-center justify-center"
        class:border-none={storeValues.isMissionSuccessful(i+1) != null}
        class:text-gray-900={storeValues.isMissionSuccessful(i+1) != null || storeValues.isCurrentMission(i+1)}
        class:bg-blue-400={storeValues.isCurrentMission(i+1)}
        class:bg-green-400={storeValues.isMissionSuccessful(i+1) === true}
        class:bg-red-400={storeValues.isMissionSuccessful(i+1) === false}
      >
      {#if storeValues.isMissionSuccessful(i+1) === null}
        {requirement.nbPeopleOnMission}
      {:else if storeValues.isMissionSuccessful(i+1) === true}
        <span class="text-5xl">&#x2713;</span>
      {:else if storeValues.isMissionSuccessful(i+1) === false}
        <span class="text-5xl">&#x2715;</span>
      {/if}
      </div>
    {/each}
  </div>
  {#if storeValues.currentGamePhase === GamePhase.TeamSelection}
    <TeamSelection dispatcher={dispatcher} storeValues={storeValues}/>
  {:else if storeValues.currentGamePhase === GamePhase.TeamVote}
    <TeamVote dispatcher={dispatcher} storeValues={storeValues}/>
  {:else if storeValues.currentGamePhase === GamePhase.Mission}
    <Mission dispatcher={dispatcher} storeValues={storeValues}/>
  {/if}
</div>
{#if storeValues.dialogShown == DialogValues.Identity}
  <Dialog dispatcher={dispatcher}>
    <Identity storeValues={storeValues}/>
  </Dialog>
{/if}