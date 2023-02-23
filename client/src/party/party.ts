import type { HttpPost } from "../http/post";
import { CreateParty, JoinParty } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";
import { CreatePartySucceeded, JoinPartySucceeded } from "../messages/events";
import type { Message } from "../messages/message-bus";

export class Party {
  constructor(
    private readonly http: HttpPost,
    private readonly dispatcher: Dispatcher,
  ){}
  
  consume(message: Message): void {
    if(message instanceof CreateParty) {
      this.http.post('/party/create', {name: message.name}).then(
        () => this.dispatcher.dispatch(new CreatePartySucceeded()),
      );
    }
    if(message instanceof JoinParty) {
      this.http.post('/party/join', {name: message.name, code: message.code}).then(
        () => this.dispatcher.dispatch(new JoinPartySucceeded()),
      );
    }
  }
}