import { FailMission, SucceedMission } from "../../messages/commands";
import type { Dispatcher } from "../../messages/dispatcher";

export interface MissionConductingValues {
  readonly player: string,
  readonly currentTeam: Set<string>,
  readonly peopleThatWorkedOnMission: Set<string>,
  readonly playerMissionSuccess: boolean,
}

export class MissionConductingService {
  constructor(
    readonly values: MissionConductingValues,
    readonly dispatcher: Dispatcher) {}

  get currentTeam(): string[] {
    return Array.from(this.values.currentTeam);
  }

  get currentTeamAsString(): string {
    return this.currentTeam.slice(0,-1).join(", ") + " and " + this.currentTeam.slice(-1);
  }

  hasGivenPlayerWorkedOnMission(player: string): boolean {
    return this.values.peopleThatWorkedOnMission.has(player);
  }

  get isPlayerInCurrentMission(): boolean {
    return this.values.currentTeam.has(this.values.player);
  }

  get hasPlayerWorkedOnMission(): boolean {
    return this.values.peopleThatWorkedOnMission.has(this.values.player);
  }

  get playerMissionSuccess(): boolean {
    return this.values.playerMissionSuccess;
  }

  succeedMission() {
    this.dispatcher.dispatch(new SucceedMission());
  }

  failMission() {
    this.dispatcher.dispatch(new FailMission());
  }
}