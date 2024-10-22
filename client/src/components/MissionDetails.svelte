<script lang="ts">
import { MissionDetailsService, MissionDetailsValues } from "./MissionDetails-service";

export let missionDetailsValues: MissionDetailsValues;

const service = new MissionDetailsService(missionDetailsValues);
</script>

<div class="text-5xl">
  Mission {service.mission} Details
</div>