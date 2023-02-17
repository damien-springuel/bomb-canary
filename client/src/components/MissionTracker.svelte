<script lang="ts">
import { MissionTrackerService, type MissionTrackerValues } from "./mission-tracker-service";
export let missionTrackerValues: MissionTrackerValues;
$: service = new MissionTrackerService(missionTrackerValues);
</script>

<div class="text-5xl">
  Mission Track
</div>
<div class="flex flex-row justify-around w-full text-3xl">
  {#each service.missions as m}
    <button 
      class="bc-button bc-button-empty h-16 w-16 flex items-center justify-center"
      class:bc-button-blue={service.isCurrentMission(m)}
      class:bc-button-green={service.shouldMissionTagShowSuccess(m)}
      class:bc-button-red={service.shouldMissionTagShowFailure(m)}
    >
      {#if service.shouldMissionTagShowNbOfPeopleOnMission(m)}
        {service.getNumberPeopleOnMission(m)} {#if service.doesMissionNeedMoreThanOneFail(m)}*{/if}
      {:else if service.shouldMissionTagShowSuccess(m)}
        <span class="text-5xl text-gray-900">&#x2713;</span>
      {:else if service.shouldMissionTagShowFailure(m)}
        <span class="text-5xl text-gray-900">&#x2715;</span>
      {/if}
    </button>
  {/each}
</div>