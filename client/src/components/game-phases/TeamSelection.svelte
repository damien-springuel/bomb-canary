<script lang="ts">
import type { Dispatcher } from "../../messages/dispatcher";
import { MissionTrackerService, type MissionTrackerValues } from "../MissionTracker-service";
import { TeamSelectionService, type TeamSelectionValues } from "./TeamSelection-service";
export let teamSelectionValues: TeamSelectionValues;
export let missionTrackerValues: MissionTrackerValues;
export let dispatcher: Dispatcher;

$: service = new TeamSelectionService(teamSelectionValues, new MissionTrackerService(missionTrackerValues), dispatcher);
</script>

<div class="flex flex-col items-center h-full w-full">
  <div class="text-xl text-center">
    <span class="font-bold">{service.leader}</span> is choosing current team.
    <div>Tentative {service.currentTeamVoteNb} of 5</div>
  </div>
  <div class="flex-grow grid grid-cols-2 w-full content-start text-center gap-2 mt-2">
    {#if service.isPlayerTheLeader}
      {#each service.players as player}
        <button 
          class="bc-button bc-button-blue" 
          class:bc-button-green={service.isGivenPlayerInTeam(player)} 
          on:click={() => service.togglePlayerSelection(player)}
          disabled={!service.isGivenPlayerSelectableForTeam(player)}
        >
          {player}
        </button>
      {/each}
    {:else}
      {#each service.players as player}
        <div 
          class="bc-tag" 
          class:bc-tag-solid={service.isGivenPlayerInTeam(player)} 
        >
          {player}
      </div>
      {/each}
    {/if}
  </div>
  {#if service.isPlayerTheLeader}
    <button 
      class="bc-button bc-button-blue" 
      on:click={() => service.confirmTeam()}
      disabled={!service.canConfirmTeam}
    >
      I'm done
    </button>
  {/if}
</div>