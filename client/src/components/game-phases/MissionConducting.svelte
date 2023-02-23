<script lang="ts">
import type { Dispatcher } from "../../messages/dispatcher";
import { MissionConductingService, type MissionConductingValues } from "./mission-conducting-service";
export let missionConductingValues: MissionConductingValues;
export let dispatcher: Dispatcher;
$: service = new MissionConductingService(missionConductingValues, dispatcher);

</script>

<div class="flex flex-col items-center h-full w-full">
  <div class="text-3xl mt-8 mb-4">
    Conducting Mission
  </div>
  <div class="grid grid-cols-2 w-full content-start gap-2 text-lg text-center">
    {#each service.currentTeam as member}
      <div 
        class="bc-tag"
        class:bc-tag-solid={service.hasGivenPlayerWorkedOnMission(member)}
      >
        {member}
      </div>
    {/each}
  </div>
  {#if service.isPlayerInCurrentMission}
    <div class="flex-grow content-center flex flex-col justify-center">
      <div class="text-4xl text-center mb-4">You are part of the mission!</div>
      {#if service.hasPlayerWorkedOnMission}
        <div class="text-center">
          You executed your part of the mission and made it
          {#if service.playerMissionSuccess}
            <span class="text-green-400">Succeed</span>
            {:else}
            <span class="text-red-400">Fail</span>
          {/if}
          .
        </div>
      {:else}
        <div class="text-2xl text-center mt-4">
          Do you want the mission to 
          <span class="text-green-400">succeed</span> 
          or to 
          <span class="text-red-400">fail</span>
          ?
        </div>
        <div class="grid grid-cols-2 justify-center w-full gap-x-2 mt-4">
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
    <div class="flex-grow content-center flex flex-col justify-center text-center text-3xl">
      {service.currentTeamAsString} are conducting the mission.
    </div>
  {/if}
</div>