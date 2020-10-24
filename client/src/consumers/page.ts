import { PartyCreated, PlayerJoined, ServerConnectionClosed } from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface RoomStore {
  showLobby: () => void,
  showPartyRoom: (code: string) => void,
}

export class PageManager {

  constructor(private readonly store: RoomStore) {}

  consume(message: Message) {
    if (message instanceof ServerConnectionClosed) {
      this.store.showLobby();
    } 
    else if(
      message instanceof PartyCreated || 
      message instanceof PlayerJoined) {
      
      this.store.showPartyRoom(message.partyCode);
    }
  }
}