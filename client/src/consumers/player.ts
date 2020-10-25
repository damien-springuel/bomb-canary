import { PlayerJoined } from "../messages/events";
import type { Message } from "../messages/messagebus";

export class PlayerManager {

  constructor(private readonly playerStore: {joinPlayer: (name: string) => void}){}

  consume(message: Message) {
    if (message instanceof PlayerJoined) {
      this.playerStore.joinPlayer(message.name);
    }
  }
}