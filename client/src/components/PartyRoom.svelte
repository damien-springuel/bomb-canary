<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { PartyRoomService, type PartyRoomValues } from "./PartyRoom-service";
export let dispatcher: Dispatcher;
export let partyRoomValues: PartyRoomValues;

$: service = new PartyRoomService(partyRoomValues, dispatcher);
let name: string;
</script>


<div class="flex flex-col items-center h-full bg-gray-900 text-blue-500 p-2 gap-5">
  <div class="text-5xl">
    Bomb Canary
  </div>
  <div class="text-2xl text-center text-blue-400">
    A real-time app to play "The Resistance" board game with friends.
  </div>
  {#if !service.hasPlayerJoined}
    <div class="flex m-2">
      <input type="text" placeholder="Name" class="bc-input m-2" bind:value={name}>
      <button class="bc-button bc-button-blue" on:click={()=>service.joinParty(name)}>Join</button>
    </div>
  {:else}
    <div>
      <div class="text-2xl underline font-bold">
        Players
      </div>
      <ul class="text-2xl">
        {#each service.players as player}
          <li>{player}</li>
        {/each}
      </ul>
    </div>
    <div class="mt-6">
      <button class="bc-button bc-button-blue" disabled={!service.canStartGame} on:click={()=>service.startGame()}>
        Start Game
      </button>
    </div>
  {/if}
</div>