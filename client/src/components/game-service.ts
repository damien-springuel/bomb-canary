export interface MissionResult {
  readonly success: boolean
  readonly nbFails: number
}

export interface MissionRequirement {
  readonly nbPeopleOnMission: number, 
  readonly nbFailuresRequiredToFail: number
}

export interface GameValues {
  missionRequirements: MissionRequirement[],
  missionResults: MissionResult[],
}

export class GameService {
  constructor(readonly gameValues: GameValues){}
  
  private get currentMission():number {
    return this.gameValues.missionResults.length;
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
    return mission < this.currentMission && this.gameValues.missionResults[mission].success;
  }

  shouldMissionTagShowFailure(mission: number): boolean {
    return mission < this.currentMission && !this.gameValues.missionResults[mission].success;
  }

  shouldMissionTagShowNbOfPeopleOnMission(mission: number): boolean {
    return mission >= this.currentMission;
  }

  getNumberPeopleOnMission(mission: number): number {
    return this.gameValues.missionRequirements[mission].nbPeopleOnMission;
  }

  get missions(): number[] {
    return [0,1,2,3,4];
  }
}