import { CreateParty, JoinParty } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";

export class LobbyService {
  constructor(private readonly dispatcher: Dispatcher){}

  joinParty(name: string, code: string) {
    this.dispatcher.dispatch(new JoinParty(name, code));
  }

  createParty(name: string) {
    this.dispatcher.dispatch(new CreateParty(name));
  }
}