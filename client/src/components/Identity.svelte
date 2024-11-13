<script lang="ts">
import { IdentityService, type IdentityValues } from './Identity-service';
export let identityValues: IdentityValues;
$: service = new IdentityService(identityValues);
</script>

<div class="bc-flex-col bc-justify-center">
  <div>You are a </div>
  {#if service.isPlayerIsASpy()}
    <div class="bc-text-title bc-text-red">Spy</div>
    <div>along with</div>
    <div class="bc-text-subtitle bc-text-red">{service.otherSpies()}</div>
  {:else}
    <div class="bc-text-title">Resistance Agent</div>
  {/if}
</div>