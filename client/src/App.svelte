<script lang="ts">
import "../tailwind.css"
import type { Store, StoreValues } from "./store/store";
import { onMount } from "svelte";
import type { Dispatcher } from "./messages/dispatcher";
import { AppService } from "./App-service";
import Page from "./components/Page.svelte";

export let dispatcher: Dispatcher;
let service = new AppService(dispatcher);

export let store: Store;

let storeValues: StoreValues;
$: storeValues = $store;

onMount(() => service.appMounted());
</script>
<Page dispatcher={dispatcher} storeValues={storeValues}/>