import { AppLoaded, CreatePartySucceeded } from "../messages/events";
import type { Message } from "../messages/messagebus";

export class Opener {
  constructor(private readonly creator: {create: ()=>void}){}

  consume(message: Message): void {
    if(message instanceof AppLoaded || message instanceof CreatePartySucceeded) {
      this.creator.create();
    }
  }
}