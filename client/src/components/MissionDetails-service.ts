import type { MissionRequirement, TeamVotes } from "../types/types";

export enum MissionTimeline {
  Past = "Past",
  Current = "Current",
  Future = "Future",
}

export interface MissionDetailsValues {
  readonly mission: number;
  readonly missionTimeline: MissionTimeline;
  readonly teamVotes: TeamVotes;
  readonly missionRequirement: MissionRequirement;
  readonly missionResult: boolean;
  readonly nbFailures: number;
}

export class MissionDetailsService{
  constructor(readonly missionDetailsValues: MissionDetailsValues){}

  get mission(): number {
    return this.missionDetailsValues.mission;
  }

  get missionTimeLine(): MissionTimeline {
    return this.missionDetailsValues.missionTimeline;
  }

  get teamVotes(): TeamVotes {
    return this.missionDetailsValues.teamVotes;
  }

  get teamSize(): number {
    return this.missionDetailsValues.missionRequirement.nbPeopleOnMission;
  }

  get nbFailuresRequiredToFail(): number {
    return this.missionDetailsValues.missionRequirement.nbFailuresRequiredToFail;
  }

  teamFromVoteAsString(vote: number): string {
    const team = Array.from(
      this.missionDetailsValues
        .teamVotes
        .votes[vote]
        .team
        .values());
    return team.join(", ");
  }

  get hasMissionSucceeded(): boolean {
    return this.missionDetailsValues.missionResult;
  }

  get nbSuccesses(): number {
    return this.teamSize - this.missionDetailsValues.nbFailures;
  }

  get nbFailures(): number {
    return this.missionDetailsValues.nbFailures;
  }

  get shouldShowVotes(): boolean {
    return this.missionTimeLine === MissionTimeline.Current || 
      this.missionTimeLine === MissionTimeline.Past;
  }

  get shouldShowMissionResult(): boolean {
    return this.missionTimeLine === MissionTimeline.Past;
  }
}