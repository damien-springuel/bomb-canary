<script lang="ts">
import { StartGame } from "../messages/commands";
import type { Message } from "../messages/messagebus";
import type { Store, StoreValues } from "../store/store";
export let dispatcher: {dispatch: (message: Message) => void};
export let store: Store;

let storeValues: StoreValues;
$: storeValues = $store;
function startGame() {
  dispatcher.dispatch(new StartGame())
}
</script>

<div class="flex flex-col items-center h-full bg-gray-900 text-blue-500 px-6">
  <h1 class="text-5xl">
    Party Code
  </h1>
  <h2 class="text-2xl text-center text-blue-400">
    {storeValues.partyCode}
  </h2>
  <button class="bc-button bc-button-blue" on:click={startGame}>Start Game</button>
</div>