import { CloseDialog, ViewIdentity, ViewMissionDetails } from "../messages/commands";
import { PartyCreated, ServerConnectionClosed, SpiesRevealed } from "../messages/events";
import type { Message } from "../messages/message-bus";

export interface RoomStore {
  showLobby(): void,
  showPartyRoom(code: string): void,
  showGameRoom(): void,
  showIdentity(): void,
  showMissionDetails(mission: number): void,
  closeDialog(): void,
}

export class PageConsumer {

  constructor(private readonly store: RoomStore) {}

  consume(message: Message): void {
    if (message instanceof ServerConnectionClosed) {
      this.store.showLobby();
    } 
    else if(message instanceof PartyCreated) {
      this.store.showPartyRoom(message.partyCode);
    }
    else if(message instanceof SpiesRevealed) {
      this.store.showGameRoom();
      this.store.showIdentity();
    }
    else if(message instanceof ViewIdentity) {
      this.store.showIdentity();
    }
    else if(message instanceof ViewMissionDetails) {
      this.store.showMissionDetails(message.mission);
    }
    else if(message instanceof CloseDialog) {
      this.store.closeDialog();
    }
  }
}