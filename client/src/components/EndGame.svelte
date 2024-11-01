<script lang="ts">
import { EndGameService, type EndGameValues } from "./EndGame-service";

export let endGameValues: EndGameValues;
const service = new EndGameService(endGameValues);
</script>

<div class="flex flex-col h-full justify-center text-center gap-8 text-3xl">
  <div 
    class="text-7xl" 
  >
    Game Over
  </div>
  <div 
    class="text-5xl"
    class:text-red-500={!service.playerHasWon}
    class:text-green-500={service.playerHasWon}
  >
    {#if service.playerHasWon}
    You won!
    {:else}
    You lost...
    {/if}
  </div>
  <div>
    The
    {#if service.spiesHaveWon}
    <span class="text-red-500">Spies </span>
    {:else}
    <span class="text-green-500">Resistance Agents </span>
    {/if}
    have won!
  </div>
  <div>
    The spies were {service.spiesAsString}
  </div>
</div>