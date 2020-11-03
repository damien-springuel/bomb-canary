import type { Message } from "./messagebus";

export class CreateParty implements Message {
  constructor(readonly name: string){}
}
export class JoinParty implements Message {
  constructor(readonly name: string, readonly code: string){}
}

export class StartGame implements Message {}

export class LeaderSelectsMember implements Message {
  constructor(readonly member: string){}
}
export class LeaderDeselectsMember implements Message {
  constructor(readonly member: string){}
}
export class LeaderConfirmsTeam implements Message {}