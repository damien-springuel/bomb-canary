<script lang="ts">
import type { Message } from "./messages/messagebus";
import type { Store, StoreValues } from "./store/store";
import { Page } from "./store/store";
import { AppLoaded } from "./messages/events";
import { onMount } from "svelte";
import Lobby from "./components/Lobby.svelte";
import PartyRoom from "./components/PartyRoom.svelte";
import Game from "./components/Game.svelte";
export let dispatcher: {dispatch: (message: Message) => void};
export let store: Store;

let storeValues: StoreValues;
$: storeValues = $store;

onMount(() => dispatcher.dispatch(new AppLoaded()));
// onMount(() => setTimeout(() => dispatcher.dispatch(new AppLoaded()), 1000));
</script>

{#if storeValues.pageToShow == Page.Lobby}
  <Lobby dispatcher={dispatcher}/>
{:else if storeValues.pageToShow == Page.PartyRoom}
  <PartyRoom dispatcher={dispatcher} storeValues={storeValues}/>
{:else if storeValues.pageToShow == Page.Game}
  <Game storeValues={storeValues}/>
{:else}
  Bomb canary loading
{/if}