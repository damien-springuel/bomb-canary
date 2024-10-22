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
}