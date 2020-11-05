import { GameStarted, LeaderDeselectedMember, LeaderSelectedMember, LeaderStartedToSelectMembers, MissionRequirement } from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface GameStore {
  setMissionRequirements(requirements: MissionRequirement[]): void
  assignLeader(leader: string): void
  selectPlayer(player: string): void
  deselectPlayer(player: string): void
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
  }
}