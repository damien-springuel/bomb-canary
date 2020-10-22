import type { Message } from "./messagebus";

export class CreateParty implements Message {
  constructor(readonly name:string){}
}

export class StartGame implements Message {}