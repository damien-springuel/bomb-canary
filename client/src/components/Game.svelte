<script lang="ts">
import type { StoreValues } from "../store/store";
import MissionConducting from "./MissionConducting.svelte";
import TeamSelection from "./TeamSelection.svelte";
import TeamVote from "./TeamVote.svelte";
import Identity from './Identity.svelte';
import Dialog from './Dialog.svelte';
import type { Dispatcher } from "../messages/dispatcher";
import MissionTracker from "./MissionTracker.svelte";
import { GameService } from "./Game-service";

export let storeValues: StoreValues;
export let dispatcher: Dispatcher;

$: service = new GameService(storeValues, dispatcher);
</script>

<div class="flex flex-col items-center h-full bg-gray-900 p-6 text-blue-500 space-y-4 text-2xl">
  <button class="bc-button bc-button-blue" on:click={()=> service.viewIdentity()}>Identity</button>
  <MissionTracker missionTrackerValues={storeValues}/>
  {#if service.isTeamSelectionPhase}
    <TeamSelection dispatcher={dispatcher} teamSelectionValues={storeValues} missionTrackerValues={storeValues}/>
  {:else if service.isTeamVotePhase}
    <TeamVote dispatcher={dispatcher} teamVoteValues={storeValues}/>
  {:else if service.isMissionConductingPhase}
    <MissionConducting dispatcher={dispatcher} missionConductingValues={storeValues}/>
  {/if}
</div>
{#if service.isDialogShownIdentity}
  <Dialog dispatcher={dispatcher}>
    <Identity identityValues={storeValues}/>
  </Dialog>
{/if}