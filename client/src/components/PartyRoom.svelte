<script lang="ts">
import { StartGame } from "../messages/commands";
import type { Message } from "../messages/messagebus";
import type { StoreValues } from "../store/store";
export let dispatcher: {dispatch: (message: Message) => void};
export let storeValues: StoreValues;

function startGame() {
  dispatcher.dispatch(new StartGame())
}
</script>

<div class="flex flex-col items-center h-full bg-gray-900 text-blue-500 p-6">
  <h1 class="text-5xl">
    Party Code
  </h1>
  <h2 class="text-3xl text-center text-blue-400">
    {storeValues.partyCode}
  </h2>
  <div class="w-full flex-grow">
    <h3 class="text-2xl underline font-semibold">
      Players
    </h3>
    <ul class="text-2xl">
      {#each storeValues.players as player}
        <li>{player}</li>
      {/each}
    </ul>
  </div>
  <button class="bc-button bc-button-blue" on:click={startGame}>Start Game</button>
</div>