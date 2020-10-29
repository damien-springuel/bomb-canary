import { LeaderStartedToSelectMembers } from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface GameStore {
  assignLeader(leader: string): void
}

export class GameManager {

  constructor(private readonly gameStore: GameStore) {}

  consume(message: Message): void {
    if (message instanceof LeaderStartedToSelectMembers) {
      this.gameStore.assignLeader(message.leader);
    }
  }
}