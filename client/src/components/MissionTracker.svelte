<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { MissionTrackerService, type MissionTrackerValues } from "./MissionTracker-service";
export let missionTrackerValues: MissionTrackerValues;
export let dispatcher: Dispatcher;
$: service = new MissionTrackerService(missionTrackerValues, dispatcher);
</script>

<div class="bc-flex-col">
  <div class="bc-text-title">
    Mission Track
  </div>
  <div class="bc-grid bc-grid-cols-5">
    {#each service.missions as m}
      <button 
        class="bc-button bc-button-empty"
        class:bc-button-blue={service.isCurrentMission(m)}
        class:bc-button-green={service.shouldMissionTagShowSuccess(m)}
        class:bc-button-red={service.shouldMissionTagShowFailure(m)}
        on:click={()=>service.viewMissionDetails(m)}
      >
        {#if service.shouldMissionTagShowNbOfPeopleOnMission(m)}
          <span class="bc-text-icon">{service.getNumberPeopleOnMission(m)} {#if service.doesMissionNeedMoreThanOneFail(m)}*{/if}</span>
        {:else if service.shouldMissionTagShowSuccess(m)}
          <span class="bc-text-icon">&#x2713;</span>
        {:else if service.shouldMissionTagShowFailure(m)}
          <span class="bc-text-icon">&#x2715;</span>
        {/if}
      </button>
    {/each}
  </div>
  <div class="bc-line"></div>
</div>