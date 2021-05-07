<script lang="ts">
  import { fly } from 'svelte/transition';
  import { cubicIn, cubicOut } from 'svelte/easing';
  import type { StoreValues } from "../store/store";
  import type { Dispatcher } from '../messages/dispatcher';
  import { IdentityService } from './identity-service';
  export let storeValues: StoreValues;
  export let dispatcher: Dispatcher;
  $: service = new IdentityService(dispatcher, storeValues);
  </script>

<div 
  in:fly={{duration: 225, y: document.body.clientHeight, easing: cubicOut, opacity: 1}} 
  out:fly={{duration: 195, y: document.body.clientHeight, easing: cubicIn, opacity: 1}} 
  class="fixed inset-0 z-10 bg-gray-800 text-blue-500 text-5xl">
  
  <div class="absolute top-0 right-0 mt-3 mr-3 bc-button bc-button-blue text-2xl" on:click={() => service.closeIdentity()}><span>&#x2715;</span></div>
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
</div>