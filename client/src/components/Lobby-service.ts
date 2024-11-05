import { JoinParty } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";

export class LobbyService {
  constructor(private readonly dispatcher: Dispatcher){}

  joinParty(name: string) {
    this.dispatcher.dispatch(new JoinParty(name));
  }
}