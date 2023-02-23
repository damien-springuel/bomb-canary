import { ApproveTeam, RejectTeam } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";

export interface TeamVoteValues {
  readonly player: string,
  readonly players: string[],
  readonly currentTeam: Set<string>,
  readonly peopleThatVotedOnTeam: Set<string>,
  readonly playerVote: boolean,
}

export class TeamVoteService {
  constructor(
    private readonly values: TeamVoteValues, 
    private readonly dispatcher: Dispatcher){}

  get currentTeamAsString(): string {
    const currentTeam = Array.from(this.values.currentTeam);
    return currentTeam.slice(0,-1).join(", ") + " and " + currentTeam.slice(-1);
  }

  hasGivenPlayerVoted(player: string): boolean {
    return this.values.peopleThatVotedOnTeam.has(player);
  }

  get hasCurrentPlayerVoted(): boolean {
    return this.hasGivenPlayerVoted(this.values.player);
  }

  get playerVote(): boolean {
    return this.values.playerVote;
  }

  get players(): string[] {
    return this.values.players;
  }

  approveTeam() {
    this.dispatcher.dispatch(new ApproveTeam());
  }

  rejectTeam() {
    this.dispatcher.dispatch(new RejectTeam());
  }
}