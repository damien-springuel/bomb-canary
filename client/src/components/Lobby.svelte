<script lang="ts">
import type { Message } from "../messages/messagebus";
import { LobbyService } from "./Lobby-service";
export let dispatcher: {dispatch: (message: Message) => void};

const service = new LobbyService(dispatcher);
let name: string;
let code: string;
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
    <input type="text" placeholder="Party Code" class="bc-input" bind:value={code}>
    <button class="bc-button bc-button-blue" on:click={()=>service.joinParty(name, code)}>Join Party</button>
    <div class="text-center">or</div>
    <button class="bc-button bc-button-blue" on:click={()=>service.createParty(name)}>Create Party</button>
  </div>
</div>