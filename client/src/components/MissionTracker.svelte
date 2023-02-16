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
    <div 
      class="rounded-full border border-blue-400 h-16 w-16 flex items-center justify-center"
      class:border-none={service.shouldMissionTagHaveNoBorder(m)}
      class:text-gray-900={service.shouldMissionTagTextBeGray(m)}
      class:bg-blue-400={service.isCurrentMission(m)}
      class:bg-green-400={service.shouldMissionTagShowSuccess(m)}
      class:bg-red-400={service.shouldMissionTagShowFailure(m)}
    >
    {#if service.shouldMissionTagShowNbOfPeopleOnMission(m)}
      {service.getNumberPeopleOnMission(m)} {#if service.doesMissionNeedMoreThanOneFail(m)}*{/if}
    {:else if service.shouldMissionTagShowSuccess(m)}
      <span class="text-5xl">&#x2713;</span>
    {:else if service.shouldMissionTagShowFailure(m)}
      <span class="text-5xl">&#x2715;</span>
    {/if}
    </div>
  {/each}
</div>