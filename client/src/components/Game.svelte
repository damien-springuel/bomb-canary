<script lang="ts">
import MissionConducting from "./MissionConducting.svelte";
import TeamSelection from "./TeamSelection.svelte";
import TeamVote from "./TeamVote.svelte";
import Identity from './Identity.svelte';
import Dialog from './Dialog.svelte';
import type { Dispatcher } from "../messages/dispatcher";
import MissionTracker from "./MissionTracker.svelte";
import { GameService, type GameValues } from "./Game-service";

export let gameValues: GameValues;
export let dispatcher: Dispatcher;

$: service = new GameService(gameValues, dispatcher);
</script>

<div class="flex flex-col items-center h-full bg-gray-900 p-6 text-blue-500 space-y-4 text-2xl">
  <button class="bc-button bc-button-blue" on:click={()=> service.viewIdentity()}>Identity</button>
  <MissionTracker missionTrackerValues={gameValues.missionTrackerValues}/>
  {#if service.isTeamSelectionPhase}
    <TeamSelection dispatcher={dispatcher} teamSelectionValues={gameValues.teamSelectionValues} missionTrackerValues={gameValues.missionTrackerValues}/>
  {:else if service.isTeamVotePhase}
    <TeamVote dispatcher={dispatcher} teamVoteValues={gameValues.teamVoteValues}/>
  {:else if service.isMissionConductingPhase}
    <MissionConducting dispatcher={dispatcher} missionConductingValues={gameValues.missionConductingValues}/>
  {/if}
</div>
{#if service.isDialogShownIdentity}
  <Dialog dispatcher={dispatcher}>
    <Identity identityValues={gameValues.identityValues}/>
  </Dialog>
{/if}