<script lang="ts">
import { MissionDetailsService, type MissionDetailsValues } from "./MissionDetails-service";

export let missionDetailsValues: MissionDetailsValues;

$: service = new MissionDetailsService(missionDetailsValues);
</script>

<div class="bc-flex-col">
  <div class="bc-text-subtitle">
    Mission #{service.mission + 1} Details
  </div>
  <div>
    <div class="bc-text-emphasis">Mission Requirements</div>
    <div class="bc-text-details">
      People needed on team: <span class="bc-text-emphasis">{service.teamSize}</span>
    </div>
    <div class="bc-text-details">
      Failures needed to fail: <span class="bc-text-emphasis">{service.nbFailuresRequiredToFail}</span>
    </div>
  </div>
  {#if service.shouldShowVotes}
    <div>
      <div class="bc-text-emphasis">
        Votes
      </div>
      <div class="bc-grid bc-grid-cols-3 bc-text-details">
        {#each service.teamVotes.votes as teamVote, i}
        <div>
          <div>
            #{i+1}: <span>{service.teamFromVoteAsString(i)}</span>
          </div>
          <div class="bc-tag">
            <div class="bc-text-subtitle">
              {#if teamVote.approved}
                <span class="bc-text-green">&#x2713;</span>
              {/if}
              {#if !teamVote.approved}
                <span class="bc-text-red">&#x2715;</span>
              {/if}
            </div>
            <div class="bc-grid bc-grid-cols-2 bc-no-gap">
              {#each teamVote.playerVotes.entries() as [player, vote]}
                <div 
                  class:bc-text-green={vote}
                  class:bc-text-red={!vote}
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
    <div>
      <div>
        <span class="bc-text-emphasis">Mission Result: </span>
        {#if service.hasMissionSucceeded}
          <span class=bc-text-green>Success</span>
        {:else}
          <span class=bc-text-red>Failure</span>
        {/if}
      </div>
      <div class="bc-grid bc-grid-cols-2">
        <div>
          <div>Successes</div>
          <div class="bc-line"></div>
          {#each Array(service.nbSuccesses) as i}
            <span class="bc-text-subtitle bc-text-green">&#x2713;</span>
          {/each}
        </div>
        <div>
          <div>Failures</div>
          <div class="bc-line"></div>
          {#each Array(service.nbFailures) as i}
            <span class="bc-text-subtitle bc-text-red">&#x2715;</span>
          {/each}
        </div>
      </div>
    </div>
  {/if}
</div>