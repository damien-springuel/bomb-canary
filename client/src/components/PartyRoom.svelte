<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import { PartyRoomService, type PartyRoomValues } from "./PartyRoom-service";
export let dispatcher: Dispatcher;
export let partyRoomValues: PartyRoomValues;

$: service = new PartyRoomService(partyRoomValues, dispatcher);
let name: string;
</script>


<div class="bc-flex-col">
  <div class="bc-text-title">
    Bomb Canary
  </div>
  <div class="">
    A real-time app to play "The Resistance" board game with friends.
  </div>
  {#if !service.hasPlayerJoined}
    <div>
      <input type="text" placeholder="Name" class="bc-input" bind:value={name}>
      <button class="bc-button bc-button-blue" on:click={()=>service.joinParty(name)}>Join</button>
    </div>
  {:else}
    <div>
      <div class="bc-font-emphasis">
        Players
      </div>
      <div class="bc-line"></div>
      <div>
        {#each service.players as player}
          <div>{player}</div>
        {/each}
      </div>
    </div>
    <div>
      <button class="bc-button bc-button-blue" disabled={!service.canStartGame} on:click={()=>service.startGame()}>
        Start Game
      </button>
    </div>
  {/if}
</div>