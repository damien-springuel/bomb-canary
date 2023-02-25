<script lang="ts">
import type { Dispatcher } from "../messages/dispatcher";
import Game from "./Game.svelte";
import Lobby from "./Lobby.svelte";
import { PageService, type PageValues } from "./Page-service";
import PartyRoom from "./PartyRoom.svelte";

export let dispatcher: Dispatcher;
export let pageValues: PageValues;

$: service = new PageService(pageValues);
</script>

{#if service.isPageLobby}
  <Lobby dispatcher={dispatcher}/>
{:else if service.isPagePartyRoom}
  <PartyRoom dispatcher={dispatcher} partyRoomValues={pageValues.partyRoomValues}/>
{:else if service.isPageGame}
  <Game dispatcher={dispatcher} gameValues={pageValues.gameValues}/>
{:else}
  Bomb canary loading
{/if}