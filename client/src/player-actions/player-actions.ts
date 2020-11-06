import type { HttpPost } from "../http/post";
import { ApproveTeam, LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember, RejectTeam, StartGame } from "../messages/commands";
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
    else if (message instanceof ApproveTeam) {
      this.http.post("/actions/approve-team");
    }
    else if (message instanceof RejectTeam) {
      this.http.post("/actions/reject-team");
    }
  }
}