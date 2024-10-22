<script lang="ts">
import { MissionDetailsService, MissionDetailsValues } from "./MissionDetails-service";

export let missionDetailsValues: MissionDetailsValues;

const service = new MissionDetailsService(missionDetailsValues);
</script>

<div class="flex flex-col items-center h-full w-full items-center">
  <div class="text-5xl my-2">
    Mission #{service.mission + 1} Details
  </div>

  <div class="text-3xl my-2">
    Votes
  </div>
  <div>
    <div class="grid grid-cols-2 gap-2 text-2xl">
      {#each service.teamVotes.votes as teamVote, i}
      <div>
        <div>
          Team #{i+1}: {service.getTeamFromVote(i)}
        </div>
        <div class="bc-tag flex flex-col">
          <div class="font-bold">
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
  <div class="text-3xl my-2">
    Mission Result
  </div>
</div>