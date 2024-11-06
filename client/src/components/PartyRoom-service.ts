import { JoinParty, StartGame } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";

export interface PartyRoomValues{
  readonly players: string[],
  readonly hasPlayerJoined: boolean,
}

export class PartyRoomService {
  constructor(
    private readonly values: PartyRoomValues,
    private readonly dispatcher: Dispatcher
  ) {}

  get players(): string[] {
    return this.values.players;
  }

  get hasPlayerJoined(): boolean {
    return this.values.hasPlayerJoined;
  }

  joinParty(name: string) {
    this.dispatcher.dispatch(new JoinParty(name));
  }

  startGame() {
    this.dispatcher.dispatch(new StartGame());
  }

  get canStartGame(): boolean {
    return this.values.players.length >= 5;
  }
}