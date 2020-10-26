import type { HttpPost } from "../http/post";
import { StartGame } from "../messages/commands";
import type { Message } from "../messages/messagebus";

export class PlayerActions {
  constructor(
    private readonly http: HttpPost,
  ){}

  consume(message: Message): void {
    if (message instanceof StartGame) {
      this.http.post("/actions/start-game");
    }
  }
}