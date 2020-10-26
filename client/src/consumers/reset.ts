import { ServerConnectionClosed } from "../messages/events";
import type { Message } from "../messages/messagebus";

export class ResetManager {

  constructor(private readonly resetter: {reset: () => void}){}

  consume(message: Message): void {
    if (message instanceof ServerConnectionClosed) {
      this.resetter.reset()
    }
  }
}