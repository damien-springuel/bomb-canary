import type { HttpPost } from "../http/post";
import { CreateParty } from "../messages/commands";
import { PartyCreated } from "../messages/events";
import type { Message } from "../messages/messagebus";

export interface CreatePartyResponse {
  code: string
}

export class Party {
  constructor(
    private readonly http: HttpPost<CreatePartyResponse>,
    private readonly dispatcher: {dispatch: (m:Message) => void},
  ){}
  
  consume(message: Message): void {
    if(message instanceof CreateParty) {
      this.http.post('/party/create', {name: message.name}).then(
        response => {
          this.dispatcher.dispatch(new PartyCreated(response.data.code));
        },
      );
    }
  }
}