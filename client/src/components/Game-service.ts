import { ViewIdentity } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";
import { Dialog, GamePhase } from "../types/types";

export interface GameValues {
  readonly currentGamePhase: GamePhase,
  readonly dialogShown: Dialog,
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

  get isDialogShownIdentity(): boolean {
    return this.values.dialogShown == Dialog.Identity;
  }
}