<script lang="ts">
import "../tailwind.css"
import type { Store } from "./store/store";
import { onMount } from "svelte";
import type { Dispatcher } from "./messages/dispatcher";
import { AppService } from "./App-service";
import Page from "./components/Page.svelte";
import { AppValuesBroker } from "./values-brokers";

export let dispatcher: Dispatcher;
let service = new AppService(dispatcher);

export let store: Store;

$: valuesBroker = new AppValuesBroker($store);

onMount(() => service.appMounted());
</script>
<Page dispatcher={dispatcher} pageValues={valuesBroker.pageValues}/>