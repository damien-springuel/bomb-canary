import type { Message } from "./message-bus";

export class JoinParty implements Message {
  constructor(readonly name: string){}
}

export class StartGame implements Message {}

export class LeaderSelectsMember implements Message {
  constructor(readonly member: string){}
}

export class LeaderDeselectsMember implements Message {
  constructor(readonly member: string){}
}

export class LeaderConfirmsTeam implements Message {}

export class ApproveTeam implements Message {}

export class RejectTeam implements Message {}

export class SucceedMission implements Message {}

export class FailMission implements Message {}

export class ViewIdentity implements Message {}

export class ViewMissionDetails implements Message {
  constructor(readonly mission:number){}
}

export class CloseDialog implements Message {}
