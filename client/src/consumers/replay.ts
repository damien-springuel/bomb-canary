import { EventsReplayEnded, ServerConnectionClosed, ServerConnectionErrorOccured } from "../messages/events";
import type { Message } from "../messages/messagebus";

export class ReplayManager {
  constructor(private readonly replayEnder: {endReplay: () => void}){}

  consume(message: Message) {
    if (message instanceof EventsReplayEnded ||
      message instanceof ServerConnectionClosed ||
      message instanceof ServerConnectionErrorOccured) {
      this.replayEnder.endReplay();
    }
  }
}
