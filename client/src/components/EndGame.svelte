<script lang="ts">
import { EndGameService, type EndGameValues } from "./EndGame-service";

export let endGameValues: EndGameValues;
$: service = new EndGameService(endGameValues);
</script>

<div class="bc-flex-col">
  <div class="bc-text-subtitle">
    {#if service.playerHasWon}
    You won!
    {:else}
    You lost...
    {/if}
  </div>
  <div>
    The
    {#if service.spiesHaveWon}
    <span class="bc-text-red">Spies </span>
    {:else}
    <span class="bc-text-green">Resistance Agents </span>
    {/if}
    have won!
  </div>
  <div>
    The spies were {service.spiesAsString}.
  </div>
</div>