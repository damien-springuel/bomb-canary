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
import { GameService } from "./game-service";

export let storeValues: StoreValues;
export let dispatcher: Dispatcher;
$: service = new GameService(storeValues);

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
    {#each service.missions as m}
      <div 
        class="rounded-full border border-blue-400 h-16 w-16 flex items-center justify-center"
        class:border-none={service.shouldMissionTagHaveNoBorder(m)}
        class:text-gray-900={service.shouldMissionTagTextBeGray(m)}
        class:bg-blue-400={service.isCurrentMission(m)}
        class:bg-green-400={service.shouldMissionTagShowSuccess(m)}
        class:bg-red-400={service.shouldMissionTagShowFailure(m)}
      >
      {#if service.shouldMissionTagShowNbOfPeopleOnMission(m)}
        {service.getNumberPeopleOnMission(m)}
      {:else if service.shouldMissionTagShowSuccess(m)}
        <span class="text-5xl">&#x2713;</span>
      {:else if service.shouldMissionTagShowFailure(m)}
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