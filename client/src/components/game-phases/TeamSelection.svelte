<script lang="ts">
import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember } from "../../messages/commands";
import type { Message } from "../../messages/messagebus";
import type { StoreValues } from "../../store/store";
import { TeamSelectionService } from "./team-selection-service";
export let storeValues: StoreValues;
export let dispatcher: {dispatch: (message: Message) => void};

$: service = new TeamSelectionService(storeValues);

function togglePlayerSelection(member: string): void {
  if (service.isPlayerTheLeader) {
    if (service.isGivenPlayerInTeam(member)) {
      dispatcher.dispatch(new LeaderDeselectsMember(member));
    } else {
      dispatcher.dispatch(new LeaderSelectsMember(member));
    }
  }
}

function confirmTeam(): void {
  dispatcher.dispatch(new LeaderConfirmsTeam());
}
</script>

<div class="flex flex-col items-center h-full w-full">
  <div class="text-xl text-center">
    <span class="font-bold">{service.leader}</span> is choosing current team.
    <div>Tentative {service.currentTeamVoteNb} of 5</div>
  </div>
  <div class="flex-grow grid grid-cols-2 w-full content-start gap-2 mt-2">
    {#each service.players as player}
      <button 
        class="bc-button bc-button-blue" 
        class:bc-button-green={service.isGivenPlayerInTeam(player)} 
        on:click={() => togglePlayerSelection(player)}
        disabled={!storeValues.isPlayerSelectableForTeam(player)}
      >
        {player}
      </button>
    {/each}
  </div>
  {#if service.isPlayerTheLeader}
    <button 
      class="bc-button bc-button-blue" 
      on:click={confirmTeam}
      disabled={!storeValues.canConfirmTeam}
      class:bc-button-gray={!storeValues.canConfirmTeam}
    >
      I'm done
    </button>
  {/if}
</div>