import type { MissionRequirement, MissionResult } from "../types/types";

export interface MissionTrackerValues {
  readonly missionRequirements: MissionRequirement[],
  readonly missionResults: MissionResult[],
}

export class MissionTrackerService {
  constructor(readonly values: MissionTrackerValues){}
  
  private get currentMission():number {
    return this.values.missionResults.length;
  }

  isCurrentMission(mission: number): boolean{
    return this.currentMission == mission;
  }

  shouldMissionTagShowSuccess(mission: number): boolean {
    return mission < this.currentMission && this.values.missionResults[mission].success;
  }

  shouldMissionTagShowFailure(mission: number): boolean {
    return mission < this.currentMission && !this.values.missionResults[mission].success;
  }

  shouldMissionTagShowNbOfPeopleOnMission(mission: number): boolean {
    return mission >= this.currentMission;
  }

  getNumberPeopleOnMission(mission: number): number {
    return this.values.missionRequirements[mission].nbPeopleOnMission;
  }

  get missions(): number[] {
    return [0,1,2,3,4];
  }

  doesMissionNeedMoreThanOneFail(mission: number): boolean {
    return this.values.missionRequirements[mission].nbFailuresRequiredToFail > 1;
  }

  get nbPeopleRequiredOnCurrentMission(): number {
    return this.values.missionRequirements[this.currentMission].nbPeopleOnMission;
  }
}