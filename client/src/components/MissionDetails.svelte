<script lang="ts">
import { MissionDetailsService, type MissionDetailsValues } from "./MissionDetails-service";

export let missionDetailsValues: MissionDetailsValues;

const service = new MissionDetailsService(missionDetailsValues);
</script>

<div class="flex flex-col items-center h-full gap-y-8">
  <div class="text-5xl my-2 underline">
    Mission #{service.mission + 1} Details
  </div>
  <div class="text-center text-xl">
    <div class="text-3xl">Mission Requirements</div>
    <div>The mission needs a team of <span class="font-bold">{service.teamSize}</span> people.</div>
    <div>
      The mission needs at least
      <span class="font-bold">{service.nbFailuresRequiredToFail}</span>
      <span class="text-red-500">failure{#if service.nbFailuresRequiredToFail > 1}<span>s</span>{/if}</span>
      to fail.
    </div>
  </div>
  {#if service.shouldShowVotes}
    <div>
      <div class="text-3xl text-center">
        Votes
      </div>
      <div class="grid grid-cols-2 gap-2 text-2xl">
        {#each service.teamVotes.votes as teamVote, i}
        <div>
          <div>
            Team #{i+1}: <span class="font-bold">{service.teamFromVoteAsString(i)}</span>
          </div>
          <div class="bc-tag flex flex-col">
            <div class="font-bold text-center underline">
              {#if teamVote.approved}
                <span class="text-green-500">Approved</span>
              {/if}
              {#if !teamVote.approved}
                <span class="text-red-500">Rejected</span>
              {/if}
            </div>
            <div class="grid grid-cols-2 gap-x-2">
              {#each teamVote.playerVotes.entries() as [player, vote]}
                <div 
                  class:text-green-500={vote}
                  class:text-red-500={!vote}
                >
                  {player}
                </div>
              {/each}
            </div>
          </div>
        </div>
        {/each}
      </div>
    </div>
  {/if}
  {#if service.shouldShowMissionResult}
    <div class="text-3xl my-2">
      <div>
        Mission Result: 
        {#if service.hasMissionSucceeded}
          <span class=text-green-500>Success</span>
        {:else}
          <span class=text-red-500>Failure</span>
        {/if}
      </div>
      <div class="grid grid-cols-2 gap-2">
        <div class="flex flex-col text-center">
          <span class="text-2xl underline">Successes</span>
          {#each Array(service.nbSuccesses) as i}
            <span class="text-5xl text-green-500">&#x2713;</span>
          {/each}
        </div>
        <div class="flex flex-col text-center">
          <span class="text-2xl underline">Failures</span>
          {#each Array(service.nbFailures) as i}
            <span class="text-5xl text-red-500">&#x2715;</span>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>