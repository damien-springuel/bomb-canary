<script lang="ts">
import { LastMissionResultService, type LastMissionResultValues } from "./LastMissionResult-service";

export let lastMissionResultValues: LastMissionResultValues;

$: service = new LastMissionResultService(lastMissionResultValues);
</script>
  
<div class="flex flex-col h-full items-center justify-center gap-y-8">
  <div class="text-4xl">The last mission was a</div>
  {#if service.success}
    <span class="text-green-500">Success</span>
    {:else}
    <span class="text-red-500">Failure</span>
  {/if}

</div>