import { ServerConnectionClosed } from "../messages/events";
import type { Message } from "../messages/messagebus";

export class PageManager {

  constructor(
    private readonly store: {
      showLobby: () => void,
    }
  ) {}

  consume(message: Message) {
    if (message instanceof ServerConnectionClosed) {
      this.store.showLobby();
    }
  }
}