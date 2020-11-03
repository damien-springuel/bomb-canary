<script lang="ts">
import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember } from "../messages/commands";

import type { Message } from "../messages/messagebus";
import type { StoreValues } from "../store/store";
export let storeValues: StoreValues;
export let dispatcher: {dispatch: (message: Message) => void};

function togglePlayerSelection(member: string): void {
  if (storeValues.isPlayerTheLeader) {
    if (storeValues.isPlayerInTeam(member)) {
      dispatcher.dispatch(new LeaderDeselectsMember(member));
    } else {
      dispatcher.dispatch(new LeaderSelectsMember(member));
    }
  }
}

function confirmTeam(): void {
  dispatcher.dispatch(new LeaderConfirmsTeam());
}
</script>

<div class="flex flex-col items-center h-full bg-gray-900 p-6 text-blue-500 space-y-4 text-2xl">
  <h2 class="text-5xl">
    Mission Track
  </h2>
  <div class="flex flex-row justify-around w-full text-gray-900 text-3xl">
    <div class="rounded-full bg-blue-400 h-16 w-16 flex items-center justify-center">1</div>
    <div class="rounded-full bg-blue-400 h-16 w-16 flex items-center justify-center">2</div>
    <div class="rounded-full bg-blue-400 h-16 w-16 flex items-center justify-center">3</div>
    <div class="rounded-full bg-blue-400 h-16 w-16 flex items-center justify-center">4</div>
    <div class="rounded-full bg-blue-400 h-16 w-16 flex items-center justify-center">5</div>
  </div>
  <div class="text-xl">
    <span class="font-bold">{storeValues.leader}</span> is choosing current team.
  </div>
  <div class="flex-grow grid grid-cols-2 w-full content-start gap-2">
    {#each storeValues.players as player}
      <button class="bc-button bc-button-blue" on:click={() => togglePlayerSelection(player)}>
        {player} (in? {storeValues.isPlayerInTeam(player)})
      </button>
    {/each}
  </div>
  {#if storeValues.isPlayerTheLeader}
    <button class="bc-button bc-button-blue" on:click={confirmTeam}>I'm done</button>
  {/if}
</div>