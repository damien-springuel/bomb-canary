import type { Message } from "./messagebus";

export class PartyCreated implements Message{
  constructor(readonly partyCode: string){}
}