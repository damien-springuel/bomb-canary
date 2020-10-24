import { EventsReplayEnded, PartyCreated, PlayerConnected, PlayerDisconnected, PlayerJoined, ServerConnectionClosed, ServerConnectionErrorOccured } from "../messages/events";
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
    if (event.EventsReplayEnded) {
      this.dispatcher.dispatch(new EventsReplayEnded());
    }
    else if (event.PartyCreated) {
      this.dispatcher.dispatch(new PartyCreated(event.PartyCreated.Code));
    }
    else if (event.PlayerConnected) {
      this.dispatcher.dispatch(new PlayerConnected(event.PlayerConnected.Name));
    }
    else if (event.PlayerDisconnected) {
      this.dispatcher.dispatch(new PlayerDisconnected(event.PlayerDisconnected.Name));
    }
    else if (event.PlayerJoined) {
      this.dispatcher.dispatch(new PlayerJoined(event.PlayerJoined.Name));
    }
  }
}