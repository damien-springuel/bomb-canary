export interface MissionResult {
  readonly success: boolean
  readonly nbFails: number
}

export interface MissionRequirement {
  readonly nbPeopleOnMission: number, 
  readonly nbFailuresRequiredToFail: number
}

export interface MissionTrackerValues {
  missionRequirements: MissionRequirement[],
  missionResults: MissionResult[],
}

export class MissionTrackerService {
  constructor(readonly values: MissionTrackerValues){}
  
  private get currentMission():number {
    return this.values.missionResults.length;
  }

  isCurrentMission(mission: number): boolean{
    return this.currentMission == mission;
  }

  shouldMissionTagHaveNoBorder(mission: number): boolean {
    return mission <= this.currentMission;
  }

  shouldMissionTagTextBeGray(mission: number): boolean {
    return mission <= this.currentMission;
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
}