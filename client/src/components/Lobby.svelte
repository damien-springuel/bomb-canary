<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { LobbyService } from "./Lobby-service";
export let dispatcher: Dispatcher;

const service = new LobbyService(dispatcher);
let name: string;
</script>

<div class="flex flex-col items-center h-full bg-gray-900 text-blue-500 px-6 overflow-auto">
  <h1 class="text-5xl">
    Bomb Canary
  </h1>
  <h2 class="text-2xl text-center text-blue-400">
    A real-time app to play "The Resistance" board game with friends.
  </h2>

  <div class="flex-grow flex flex-col justify-evenly w-full text-2xl items-stretch space-y-10 mb-4">
    <input type="text" placeholder="Name" class="bc-input" bind:value={name}>
    <div class="border-solid border-t-2 border-blue-500"></div>
    <button class="bc-button bc-button-blue" on:click={()=>service.joinParty(name)}>Join Party</button>
  </div>
</div>