import { AppLoaded, PartyCreated } from "../messages/events";
import type { Message } from "../messages/messagebus";

export class Opener {
  constructor(private readonly opener: {open: ()=>void}){}

  consume(message: Message): void {
    if(message instanceof AppLoaded || message instanceof PartyCreated) {
      this.opener.open();
    }
  }
}