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
  class="bc-dialog"
>
  <button 
    class="bc-absolute-top-right bc-button bc-button-blue" 
    on:click={()=> service.closeDialog()}
  >
    <span class="bc-text-icon">&#x2715;</span>
  </button>
  <slot></slot>
</div>