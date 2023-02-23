<script lang="ts">
import { fly } from 'svelte/transition';
import { cubicIn, cubicOut } from 'svelte/easing';
import type { Dispatcher } from '../messages/dispatcher';
import { DialogService } from './Dialog-service';
export let dispatcher: Dispatcher;

$: service = new DialogService(dispatcher);
</script>

<div 
  in:fly={{duration: 225, y: document.body.clientHeight, easing: cubicOut, opacity: 1}} 
  out:fly={{duration: 195, y: document.body.clientHeight, easing: cubicIn, opacity: 1}} 
  class="fixed inset-0 z-10 bg-gray-800 text-blue-500 text-5xl"
>
  <button 
    class="absolute top-0 right-0 mt-3 mr-3 bc-button bc-button-blue h-16 w-16 flex items-center justify-center" 
    on:click={()=> service.closeDialog()}
  >
    <span class="text-4xl">&#x2715;</span>
  </button>
  <slot></slot>
</div>