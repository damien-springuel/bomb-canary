<script lang="ts">
import MissionConducting from "./MissionConducting.svelte";
import TeamSelection from "./TeamSelection.svelte";
import TeamVote from "./TeamVote.svelte";
import Identity from './Identity.svelte';
import Dialog from './Dialog.svelte';
import type { Dispatcher } from "../messages/dispatcher";
import MissionTracker from "./MissionTracker.svelte";
import { GameService, type GameValues } from "./Game-service";
import MissionDetails from "./MissionDetails.svelte";
import LastMissionResult from "./LastMissionResult.svelte";
import EndGame from "./EndGame.svelte";

export let gameValues: GameValues;
export let dispatcher: Dispatcher;

$: service = new GameService(gameValues, dispatcher);
</script>

<div class="bc-flex-col">
  <button class="bc-button bc-button-blue" on:click={()=> service.viewIdentity()}>Identity</button>
  <MissionTracker missionTrackerValues={gameValues.missionTrackerValues} dispatcher={dispatcher}/>
  {#if service.isTeamSelectionPhase}
    <TeamSelection dispatcher={dispatcher} teamSelectionValues={gameValues.teamSelectionValues}/>
  {:else if service.isTeamVotePhase}
    <TeamVote dispatcher={dispatcher} teamVoteValues={gameValues.teamVoteValues}/>
  {:else if service.isMissionConductingPhase}
    <MissionConducting dispatcher={dispatcher} missionConductingValues={gameValues.missionConductingValues}/>
  {:else if service.gameHasEnded}
    <EndGame endGameValues={gameValues.endGameValues}/>
  {/if}
</div>
{#if service.isDialogShownIdentity}
  <Dialog dispatcher={dispatcher}>
    <Identity identityValues={gameValues.identityValues}/>
  </Dialog>
{/if}
{#if service.isDialogShownMissionDetails}
  <Dialog dispatcher={dispatcher}>
    <MissionDetails missionDetailsValues={gameValues.missionDetailsValues}/>
  </Dialog>
{/if}
{#if service.isDialogShownLastMissionResult}
  <Dialog dispatcher={dispatcher}>
    <LastMissionResult lastMissionResultValues={gameValues.lastMissionResultValues}/>
  </Dialog>
{/if}