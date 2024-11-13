<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { TeamVoteService, type TeamVoteValues } from "./TeamVote-service";
export let teamVoteValues: TeamVoteValues;
export let dispatcher: Dispatcher;

$: service = new TeamVoteService(teamVoteValues, dispatcher);
</script>

<div class="bc-flex-col">
  <div class="bc-text-emphasis">
    Team
  </div>
  <div class="bc-tag">
    {service.currentTeamAsString}
  </div>
  {#if !service.hasCurrentPlayerVoted}
    <div class="bc-grid bc-grid-cols-2">
      <button class="bc-button bc-button-green" on:click={()=> service.approveTeam()}>
        Approve
      </button>
      <button class="bc-button bc-button-red" on:click={()=> service.rejectTeam()}>
        Reject
      </button>
    </div>
  {:else}
    {#if service.playerVote}
      <div>
        You <span class="bc-text-green">approved</span> the team.
      </div> 
    {:else}
      <div>
        You <span class="bc-text-red">rejected</span> the team.
      </div> 
    {/if}
  {/if}
  <div class="bc-line"></div>
  <div class="bc-text-emphasis">
    Who voted?
  </div>
  <div class="bc-grid bc-grid-cols-3">
    {#each service.players as player}
      <div class="bc-tag" class:bc-tag-solid={service.hasGivenPlayerVoted(player)}>
        {player}
      </div>
    {/each}
  </div>
</div>