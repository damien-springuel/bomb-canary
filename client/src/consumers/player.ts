import { EventsReplayStarted, PlayerJoined, SpiesRevealed } from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface PlayerStore {
  definePlayer(name: string): void
  joinPlayer(name: string): void
  rememberSpies(spies: Set<string>): void
  showIdentity(): void
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
    else if (message instanceof SpiesRevealed) {
      this.playerStore.rememberSpies(message.spies);
      this.playerStore.showIdentity();
    }
  }
}