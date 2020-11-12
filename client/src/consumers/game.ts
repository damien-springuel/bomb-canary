import { GameStarted, LeaderConfirmedTeam, LeaderDeselectedMember, LeaderSelectedMember, LeaderStartedToSelectMembers, MissionRequirement, MissionStarted, PlayerVotedOnTeam, PlayerWorkedOnMission } from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface GameStore {
  setMissionRequirements(requirements: MissionRequirement[]): void
  assignLeader(leader: string): void
  selectPlayer(player: string): void
  deselectPlayer(player: string): void
  startTeamVote(): void
  makePlayerVote(player: string, approval: boolean | null): void
  startMission(): void
  makePlayerWorkOnMission(player: string, success: boolean | null): void
}

export class GameManager {

  constructor(private readonly gameStore: GameStore) {}

  consume(message: Message): void {
    if (message instanceof LeaderStartedToSelectMembers) {
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
    else if(message instanceof MissionStarted) {
      this.gameStore.startMission();
    }
    else if(message instanceof PlayerWorkedOnMission) {
      this.gameStore.makePlayerWorkOnMission(message.player, message.success);
    }
  }
}