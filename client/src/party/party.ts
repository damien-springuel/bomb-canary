import type { HttpPost } from "../http/post";
import { CreateParty } from "../messages/commands";
import type { Message } from "../messages/messagebus";

export class Party {
  constructor(private readonly http: HttpPost){}
  
  consume(message: Message): void {
    if(message instanceof CreateParty) {
      this.http.post('/party/create', {name: message.name});
    }
  }
}