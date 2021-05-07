<script lang="ts">
  import type { StoreValues } from "../store/store";
  import { IdentityService } from './identity-service';
  export let storeValues: StoreValues;
  $: service = new IdentityService(storeValues);
  </script>

<div class="flex flex-col h-full items-center justify-center">
  <div class="text-3xl">You are a </div>
  {#if service.isPlayerIsASpy()}
    <div class="text-red-500 mt-4">Spy</div>
    <div class="text-3xl mt-4">along with</div>
    <div class="text-red-500 text-4xl mt-4">{service.otherSpies().join(', ')}</div>
  {:else}
    <div class="text-blue-500 mt-4">Resistance Agent</div>
  {/if}
</div>