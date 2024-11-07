<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { TeamSelectionService, type TeamSelectionValues } from "./TeamSelection-service";
export let teamSelectionValues: TeamSelectionValues;
export let dispatcher: Dispatcher;

$: service = new TeamSelectionService(teamSelectionValues, dispatcher);
</script>

<div class="bc-flex-col">
  <div class="bc-text-subtitle">
    <span class="bc-text-emphasis">{service.leader}</span> is choosing current team 
    of {service.nbPeopleRequiredOnMission} people.
    <div>Tentative {service.currentTeamVoteNb} of 5.</div>
  </div>
  <div class="bc-grid bc-grid-cols-2">
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