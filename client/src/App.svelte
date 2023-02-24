<script lang="ts">
import "../tailwind.css"
import type { Store, StoreValues } from "./store/store";
import { onMount } from "svelte";
import Lobby from "./components/Lobby.svelte";
import PartyRoom from "./components/PartyRoom.svelte";
import Game from "./components/Game.svelte";
import type { Dispatcher } from "./messages/dispatcher";
import { AppService } from "./App-service";

export let dispatcher: Dispatcher;
export let store: Store;

let storeValues: StoreValues;
$: storeValues = $store;
$: service = new AppService(storeValues, dispatcher);

onMount(() => service.appMounted());
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