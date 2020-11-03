import type { HttpPost } from "../http/post";
import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember, StartGame } from "../messages/commands";
import type { Message } from "../messages/messagebus";

export class PlayerActions {
  constructor(
    private readonly http: HttpPost,
  ){}

  consume(message: Message): void {
    if (message instanceof StartGame) {
      this.http.post("/actions/start-game");
    }
    else if (message instanceof LeaderSelectsMember) {
      this.http.post("/actions/leader-selects-member", {member: message.member});
    }
    else if (message instanceof LeaderDeselectsMember) {
      this.http.post("/actions/leader-deselects-member", {member: message.member});
    }
    else if (message instanceof LeaderConfirmsTeam) {
      this.http.post("/actions/leader-confirms-team");
    }
  }
}