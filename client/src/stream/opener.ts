import { AppLoaded, JoinPartySucceeded } from "../messages/events";
import type { Message } from "../messages/message-bus";

export class Opener {
  constructor(private readonly creator: {create: ()=>void}){}

  consume(message: Message): void {
    if(message instanceof AppLoaded || 
      message instanceof JoinPartySucceeded) {
      this.creator.create();
    }
  }
}