import { EventsReplayStarted, PlayerJoined } from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface PlayerStore {
  definePlayer(name: string): void
  joinPlayer(name: string): void
}

export class PlayerManager {

  constructor(private readonly playerStore: PlayerStore){}

  consume(message: Message) {
    if (message instanceof PlayerJoined) {
      this.playerStore.joinPlayer(message.name);
    }
    else if (message instanceof EventsReplayStarted) {
      this.playerStore.definePlayer(message.playerName);
    }
  }
}