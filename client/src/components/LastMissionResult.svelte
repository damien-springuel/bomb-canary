<script lang="ts">
import { LastMissionResultService, type LastMissionResultValues } from "./LastMissionResult-service";

export let lastMissionResultValues: LastMissionResultValues;

$: service = new LastMissionResultService(lastMissionResultValues);
</script>
  
<div class="bc-flex-col bc-justify-center bc-text-title">
  <div>The last mission was a</div>
  {#if service.success}
    <span class="bc-text-green">Success!</span>
    {:else}
    <span class="bc-text-red">Failure...</span>
  {/if}
</div>