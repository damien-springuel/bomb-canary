import type { Message } from "./messagebus";

export class AppLoaded implements Message{}

export class CreatePartySucceeded implements Message{}
export class JoinPartySucceeded implements Message{}

export class ServerConnectionClosed implements Message {}
export class ServerConnectionErrorOccured implements Message {}

export class EventsReplayStarted implements Message {
  constructor(readonly playerName: string){}
}
export class EventsReplayEnded implements Message {}

export class PartyCreated implements Message{
  constructor(readonly partyCode: string){}
}
export class PlayerConnected implements Message {
  constructor(readonly name: string){}
}

export class PlayerDisconnected implements Message {
  constructor(readonly name: string){}
}

export class PlayerJoined implements Message {
  constructor(readonly name: string) {}
}

export interface MissionRequirement {
  readonly nbPeopleOnMission: number, 
  readonly nbFailuresRequiredToFail: number
}

export class GameStarted implements Message {
  constructor(readonly requirements: MissionRequirement[]) {}
}

export class SpiesRevealed implements Message {
  constructor(readonly spies: Set<string>) {}
}

export class LeaderStartedToSelectMembers implements Message {
  constructor(readonly leader: string) {}
}

export class LeaderSelectedMember implements Message {
  constructor(readonly member: string) {}
}

export class LeaderDeselectedMember implements Message {
  constructor(readonly member: string) {}
}

export class LeaderConfirmedTeam implements Message {}

export class PlayerVotedOnTeam implements Message{
  constructor(readonly player: string, readonly approved: boolean | null){}
}

export class AllPlayerVotedOnTeam implements Message{
  constructor(readonly approved: boolean, readonly voteFailures: number, readonly playerVotes: Map<string, boolean>){}
}

export class MissionStarted implements Message {}

export class PlayerWorkedOnMission implements Message {
  constructor(readonly player: string, readonly success: boolean | null){}
}

export class MissionCompleted implements Message {
  constructor(readonly success: boolean, readonly nbFails: number){}
}
