<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { PartyRoomService, type PartyRoomValues } from "./PartyRoom-service";
export let dispatcher: Dispatcher;
export let partyRoomValues: PartyRoomValues;

$: service = new PartyRoomService(partyRoomValues, dispatcher);
</script>

<div class="flex flex-col items-center h-full bg-gray-900 text-blue-500 p-6">
  <h1 class="text-5xl">
    Party Code
  </h1>
  <h2 class="text-3xl text-center text-blue-400">
    {service.partyCode}
  </h2>
  <div class="w-full">
    <h3 class="text-2xl underline font-semibold">
      Players
    </h3>
    <ul class="text-2xl">
      {#each service.players as player}
        <li>{player}</li>
      {/each}
    </ul>
  </div>
  <div class="mt-6">
    <button class="bc-button bc-button-blue" on:click={()=>service.startGame()}>
      Start Game
    </button>
  </div>
</div>