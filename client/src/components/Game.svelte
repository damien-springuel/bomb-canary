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
import MissionTracker from "./MissionTracker.svelte";

export let storeValues: StoreValues;
export let dispatcher: Dispatcher;

function viewIdentity() {
  dispatcher.dispatch(new ViewIdentity());
}
</script>

<div class="flex flex-col items-center h-full bg-gray-900 p-6 text-blue-500 space-y-4 text-2xl">
  <button class="bc-button bc-button-blue" on:click={viewIdentity}>Identity</button>
  <MissionTracker missionTrackerValues={storeValues}/>
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
    <Identity identityValues={storeValues}/>
  </Dialog>
{/if}