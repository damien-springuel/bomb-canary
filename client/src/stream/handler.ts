import { PlayerConnected, ServerConnectionClosed, ServerConnectionErrorOccured } from "../messages/events";
import type { Message } from "../messages/messagebus";
import type { ServerEvent } from "./server-event";

export class Handler {

  constructor(private readonly dispatcher: {dispatch: (m: Message) => void}){}

  onClose(): void {
    this.dispatcher.dispatch(new ServerConnectionClosed());
  }
  onError(): void {
    this.dispatcher.dispatch(new ServerConnectionErrorOccured());
  }
  onEvent(event: ServerEvent) {
    if (event.PlayerConnected) {
      this.dispatcher.dispatch(new PlayerConnected(event.PlayerConnected.Name));
    }
  }
}