import { ViewIdentity } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";
import { Dialog, GamePhase } from "../types/types";
import type { EndGameValues } from "./EndGame-service";
import type { IdentityValues } from "./Identity-service";
import type { LastMissionResultValues } from "./LastMissionResult-service";
import type { MissionConductingValues } from "./MissionConducting-service";
import type { MissionDetailsValues } from "./MissionDetails-service";
import type { MissionTrackerValues } from "./MissionTracker-service";
import type { TeamSelectionValues } from "./TeamSelection-service";
import type { TeamVoteValues } from "./TeamVote-service";

export interface GameValues {
  readonly currentGamePhase: GamePhase,
  readonly dialogShown: Dialog,
  readonly identityValues: IdentityValues;
  readonly missionTrackerValues: MissionTrackerValues;
  readonly teamSelectionValues: TeamSelectionValues;
  readonly teamVoteValues: TeamVoteValues;
  readonly missionConductingValues: MissionConductingValues;
  readonly missionDetailsValues: MissionDetailsValues;
  readonly lastMissionResultValues: LastMissionResultValues;
  readonly endGameValues: EndGameValues;
}

export class GameService {
  constructor(
    private readonly values: GameValues,
    private readonly dispatcher: Dispatcher) {}

  viewIdentity() {
    this.dispatcher.dispatch(new ViewIdentity());
  }

  private isPhase(phase: GamePhase): boolean {
    return this.values.currentGamePhase == phase;
  }

  get isTeamSelectionPhase(): boolean {
    return this.isPhase(GamePhase.TeamSelection);
  }

  get isTeamVotePhase(): boolean {
    return this.isPhase(GamePhase.TeamVote);
  }

  get isMissionConductingPhase(): boolean {
    return this.isPhase(GamePhase.Mission);
  }

  get gameHasEnded(): boolean {
    return this.isPhase(GamePhase.GameEnded);
  }

  get isDialogShownIdentity(): boolean {
    return this.values.dialogShown == Dialog.Identity;
  }

  get isDialogShownMissionDetails(): boolean {
    return this.values.dialogShown == Dialog.MissionDetails;
  }

  get isDialogShownLastMissionResult(): boolean {
    return this.values.dialogShown == Dialog.LastMissionResult;
  }
}