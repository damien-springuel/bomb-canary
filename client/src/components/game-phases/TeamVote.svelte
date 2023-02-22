<script lang="ts">
import type { Dispatcher } from "../../messages/dispatcher";
import { TeamVoteService, type TeamVoteValues } from "./team-vote-service";
export let teamVoteValues: TeamVoteValues;
export let dispatcher: Dispatcher;

$: service = new TeamVoteService(teamVoteValues, dispatcher);
</script>

<div class="flex flex-col items-center h-full w-full">
  <div class="text-3xl">
    Team
  </div>
  <div class="flex flex-row w-full text-blue-500 justify-center bc-tag mt-4 mb-4">
    {service.currentTeamAsString}
  </div>
  {#if !service.hasCurrentPlayerVoted}
    <div class="grid grid-cols-2 justify-center w-full gap-x-2">
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
        You <span class="text-green-500">approved</span> the team.
      </div> 
    {:else}
      <div>
        You <span class="text-red-500">rejected</span> the team.
      </div> 
    {/if}
  {/if}
  
  <div class="text-3xl mt-8 mb-4">
    Votes
  </div>
  <div class="flex-grow grid grid-cols-3 w-full content-start gap-2 text-lg text-center">
    {#each service.players as player}
      <div class="bc-tag" class:bc-tag-solid={service.hasGivenPlayerVoted(player)}>
        {player}
      </div>
    {/each}
  </div>
</div>