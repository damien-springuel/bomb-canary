import { EventsReplayStarted, PlayerJoined, SpiesRevealed } from "../messages/events";
import type { Message } from "../messages/message-bus";

export interface PlayerStore {
  definePlayer(name: string): void
  joinPlayer(name: string): void
  rememberSpies(spies: Set<string>): void
}

export class PlayerConsumer {

  constructor(private readonly playerStore: PlayerStore){}

  consume(message: Message) {
    if (message instanceof PlayerJoined) {
      this.playerStore.joinPlayer(message.name);
    }
    else if (message instanceof EventsReplayStarted) {
      this.playerStore.definePlayer(message.playerName);
    }
    else if (message instanceof SpiesRevealed) {
      this.playerStore.rememberSpies(message.spies);
    }
  }
}