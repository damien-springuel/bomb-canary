<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import type { StoreValues } from "../store/store";
import Game from "./Game.svelte";
import Lobby from "./Lobby.svelte";
import { PageService } from "./page-service";
import PartyRoom from "./PartyRoom.svelte";

export let dispatcher: Dispatcher;
export let storeValues: StoreValues;

$: service = new PageService(storeValues);
</script>

{#if service.isPageLobby}
  <Lobby dispatcher={dispatcher}/>
{:else if service.isPagePartyRoom}
  <PartyRoom dispatcher={dispatcher} storeValues={storeValues}/>
{:else if service.isPageGame}
  <Game dispatcher={dispatcher} storeValues={storeValues}/>
{:else}
  Bomb canary loading
{/if}