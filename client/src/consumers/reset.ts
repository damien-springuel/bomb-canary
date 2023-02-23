import { ServerConnectionClosed } from "../messages/events";
import type { Message } from "../messages/message-bus";

export class ResetManager {

  constructor(private readonly resetter: {reset: () => void}){}

  consume(message: Message): void {
    if (message instanceof ServerConnectionClosed) {
      this.resetter.reset()
    }
  }
}