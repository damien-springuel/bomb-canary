import { AppLoaded, PartyCreated } from "../messages/events";
import type { Message } from "../messages/messagebus";

export class ServerEventConnectionOpener {
  constructor(private readonly connectionOpener: {open: ()=>void}){}

  consume(message: Message): void {
    if(message instanceof AppLoaded || message instanceof PartyCreated) {
      this.connectionOpener.open();
    }
  }
}