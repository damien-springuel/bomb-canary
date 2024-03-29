import { 
  EventsReplayEnded, 
  EventsReplayStarted, 
  ServerConnectionClosed, 
  ServerConnectionErrorOccured 
} from "../messages/events";
import type { Message } from "../messages/message-bus";

export interface ReplayStore {
  startReplay(): void,
  endReplay(): void
}

export class ReplayConsumer {
  constructor(private readonly replayStore: ReplayStore){}

  consume(message: Message) {
    if (message instanceof EventsReplayEnded ||
      message instanceof ServerConnectionClosed ||
      message instanceof ServerConnectionErrorOccured) {
      this.replayStore.endReplay();
    } 
    else if (message instanceof EventsReplayStarted) {
      this.replayStore.startReplay();
    }
  }
}
