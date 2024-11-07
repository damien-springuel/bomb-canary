<script lang="ts">
import "../tailwind.css"
import type { Store } from "./store/store";
import { onMount } from "svelte";
import type { Dispatcher } from "./messages/dispatcher";
import { AppService } from "./App-service";
import { AppValuesBroker } from "./values-brokers";
import Game from "./components/Game.svelte";
import PartyRoom from "./components/PartyRoom.svelte";

export let dispatcher: Dispatcher;
export let store: Store;

$: valuesBroker = new AppValuesBroker($store);
$: service = new AppService(dispatcher, valuesBroker);

onMount(() => service.appMounted());
</script>

<div class="bc-app">
  {#if service.isPagePartyRoom}
    <PartyRoom dispatcher={dispatcher} partyRoomValues={valuesBroker.partyRoomValues}/>
  {:else if service.isPageGame}
    <Game dispatcher={dispatcher} gameValues={valuesBroker.gameValues}/>
  {:else}
    Bomb canary loading
  {/if}

</div>