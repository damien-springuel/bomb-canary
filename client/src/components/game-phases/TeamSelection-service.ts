import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember } from "../../messages/commands";
import type { Dispatcher } from "../../messages/dispatcher";
import type { MissionTrackerService } from "../MissionTracker-service";

export interface TeamSelectionValues {
  readonly currentTeam: Set<string>,
  readonly player: string,
  readonly leader: string,
  readonly currentTeamVoteNb: number,
  readonly players: string[],
}

export class TeamSelectionService {
  constructor(
    private readonly values: TeamSelectionValues, 
    private readonly missionTrackerService: MissionTrackerService,
    private readonly dispatcher: Dispatcher){}

  isGivenPlayerInTeam(player: string): boolean {
    return this.values.currentTeam.has(player);
  }

  get isPlayerTheLeader(): boolean {
    return this.values.player === this.values.leader;
  }

  get leader(): string {
    return this.values.leader;
  }

  get currentTeamVoteNb(): number {
    return this.values.currentTeamVoteNb;
  }
  
  get players(): string[] {
    return this.values.players;
  }

  isGivenPlayerSelectableForTeam(player: string): boolean {
    return this.isGivenPlayerInTeam(player) || this.values.currentTeam.size < this.missionTrackerService.nbPeopleRequiredOnCurrentMission;
  }
  
  get canConfirmTeam(): boolean {
    return this.values.currentTeam.size === this.missionTrackerService.nbPeopleRequiredOnCurrentMission;
  }

  togglePlayerSelection(player: string) {
    if(this.isGivenPlayerInTeam(player)) {
      this.dispatcher.dispatch(new LeaderDeselectsMember(player));
    } else {
      this.dispatcher.dispatch(new LeaderSelectsMember(player));
    }
  }

  confirmTeam() {
    this.dispatcher.dispatch(new LeaderConfirmsTeam());
  }
}