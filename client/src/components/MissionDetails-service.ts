import type { TeamVotes } from "../types/types";

export class MissionDetailsValues {
  readonly mission: number;
  readonly teamVotes: TeamVotes;
}

export class MissionDetailsService{
  constructor(readonly missionDetailsValues: MissionDetailsValues){}

  get mission(): number {
    return this.missionDetailsValues.mission;
  }

  get teamVotes(): TeamVotes {
    return this.missionDetailsValues.teamVotes;
  }

  getTeamFromVote(vote: number): string {
    const team = Array.from(
      this.missionDetailsValues
        .teamVotes
        .votes[vote]
        .team
        .values());
    return team.join(", ");
  }
}