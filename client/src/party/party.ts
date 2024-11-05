import type { HttpPost } from "../http/post";
import { JoinParty } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";
import { JoinPartySucceeded } from "../messages/events";
import type { Message } from "../messages/message-bus";

export class Party {
  constructor(
    private readonly http: HttpPost,
    private readonly dispatcher: Dispatcher,
  ){}
  
  consume(message: Message): void {
    if(message instanceof JoinParty) {
      this.http.post('/party/join', {name: message.name}).then(
        () => this.dispatcher.dispatch(new JoinPartySucceeded()),
      );
    }
  }
}