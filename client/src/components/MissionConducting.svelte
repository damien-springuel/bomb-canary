<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { MissionConductingService, type MissionConductingValues } from "./MissionConducting-service";

export let missionConductingValues: MissionConductingValues;
export let dispatcher: Dispatcher;

$: service = new MissionConductingService(missionConductingValues, dispatcher);
</script>

<div class="bc-flex-col">
  
  {#if service.isPlayerInCurrentMission}
    <div class="bc-flex-col">
      <div class="bc-text-subtitle">You are part of the mission!</div>
      {#if service.hasPlayerWorkedOnMission}
        <div>
          You executed your part of the mission and made it
          {#if service.playerMissionSuccess}
            <span class="bc-text-green">Succeed</span>
            {:else}
            <span class="bc-text-red">Fail</span>
          {/if}
          .
        </div>
      {:else}
        <div>
          Do you want the mission to 
          <span class="bc-text-green">succeed</span> 
          or to 
          <span class="bc-text-red">fail</span>
          ?
        </div>
        <div class="bc-grid bc-grid-cols-2">
          <button class="bc-button bc-button-green" on:click={()=> service.succeedMission()}>
            Succeed
          </button>
          <button class="bc-button bc-button-red" on:click={()=> service.failMission()}>
            Fail
          </button>
        </div>
      {/if}
    </div>
  {:else}
    <div class="bc-text-subtitle">
      {service.currentTeamAsString} are conducting the mission.
    </div>
    <div class="bc-grid bc-grid-cols-2">
    {#each service.currentTeam as member}
      <div 
        class="bc-tag"
        class:bc-tag-solid={service.hasGivenPlayerWorkedOnMission(member)}
      >
        {member}
      </div>
    {/each}
  </div>
  {/if}
  
</div>