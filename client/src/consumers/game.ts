import { 
  AllPlayerVotedOnTeam, 
  GameStarted, 
  LeaderConfirmedTeam, 
  LeaderDeselectedMember, 
  LeaderSelectedMember, 
  LeaderStartedToSelectMembers, 
  MissionCompleted, 
  MissionStarted, 
  PlayerVotedOnTeam, 
  PlayerWorkedOnMission, 
  type MissionRequirement 
} from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface GameStore {
  setMissionRequirements(requirements: MissionRequirement[]): void
  startTeamSelection(): void
  assignLeader(leader: string): void
  selectPlayer(player: string): void
  deselectPlayer(player: string): void
  startTeamVote(): void
  makePlayerVote(player: string, approval: boolean | null): void
  saveTeamVoteResult(approved: boolean, playerVotes: Map<string, boolean>): void
  startMission(): void
  makePlayerWorkOnMission(player: string, success: boolean | null): void
  saveMissionResult(success: boolean, nbFails: number): void
}

export class GameManager {

  constructor(private readonly gameStore: GameStore) {}

  consume(message: Message): void {
    if (message instanceof LeaderStartedToSelectMembers) {
      this.gameStore.startTeamSelection();
      this.gameStore.assignLeader(message.leader);
    } 
    else if(message instanceof LeaderSelectedMember) {
      this.gameStore.selectPlayer(message.member);
    }
    else if(message instanceof LeaderDeselectedMember) {
      this.gameStore.deselectPlayer(message.member);
    }
    else if(message instanceof GameStarted) {
      this.gameStore.setMissionRequirements(message.requirements);
    }
    else if(message instanceof LeaderConfirmedTeam) {
      this.gameStore.startTeamVote();
    }
    else if(message instanceof PlayerVotedOnTeam) {
      this.gameStore.makePlayerVote(message.player, message.approved);
    }
    else if(message instanceof AllPlayerVotedOnTeam) {
      this.gameStore.saveTeamVoteResult(message.approved, message.playerVotes);
    }
    else if(message instanceof MissionStarted) {
      this.gameStore.startMission();
    }
    else if(message instanceof PlayerWorkedOnMission) {
      this.gameStore.makePlayerWorkOnMission(message.player, message.success);
    }
    else if(message instanceof MissionCompleted) {
      this.gameStore.saveMissionResult(message.success, message.nbFails);
    }
  }
}