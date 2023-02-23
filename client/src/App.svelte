<script lang="ts">
import "../tailwind.css"
import type { Store, StoreValues } from "./store/store";
import { Page } from "./store/store";
import { AppLoaded } from "./messages/events";
import { onMount } from "svelte";
import Lobby from "./components/Lobby.svelte";
import PartyRoom from "./components/PartyRoom.svelte";
import Game from "./components/Game.svelte";
import type { Dispatcher } from "./messages/dispatcher";

export let dispatcher: Dispatcher;
export let store: Store;

let storeValues: StoreValues;
$: storeValues = $store;

onMount(() => dispatcher.dispatch(new AppLoaded()));
</script>

{#if storeValues.pageToShow == Page.Lobby}
  <Lobby dispatcher={dispatcher}/>
{:else if storeValues.pageToShow == Page.PartyRoom}
  <PartyRoom dispatcher={dispatcher} storeValues={storeValues}/>
{:else if storeValues.pageToShow == Page.Game}
  <Game dispatcher={dispatcher} storeValues={storeValues}/>
{:else}
  Bomb canary loading
{/if}