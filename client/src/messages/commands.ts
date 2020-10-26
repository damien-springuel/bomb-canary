import type { Message } from "./messagebus";

export class CreateParty implements Message {
  constructor(readonly name: string){}
}
export class JoinParty implements Message {
  constructor(readonly name: string, readonly code: string){}
}

export class StartGame implements Message {}