import { StartGame } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";

export interface PartyRoomValues{
  readonly partyCode: string,
  readonly players: string[],
}

export class PartyRoomService {
  constructor(
    private readonly values: PartyRoomValues,
    private readonly dispatcher: Dispatcher
  ) {}

  get partyCode(): string {
    return this.values.partyCode;
  }

  get players(): string[] {
    return this.values.players;
  }

  startGame() {
    this.dispatcher.dispatch(new StartGame());
  }
}