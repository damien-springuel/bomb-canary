export enum Page {
  Loading = "loading",
  PartyRoom = "partyRoom",
  Game = "game",
}

export enum GamePhase {
  TeamSelection = "teamSelection",
  TeamVote = "teamVote",
  Mission = "mission",
  GameEnded = "gameEnded",
}

export enum Dialog {
  Identity = "identity",
  MissionDetails = "missionDetails",
  LastMissionResult = "LastMissionResult",
}

export interface MissionResult {
  readonly success: boolean
  readonly nbFails: number
}

export interface TeamVote {
  readonly team: Set<string>
  readonly approved: boolean
  readonly playerVotes: Map<string, boolean>
}

export interface TeamVotes {
  readonly votes: TeamVote[]
}

export interface MissionRequirement {
  readonly nbPeopleOnMission: number, 
  readonly nbFailuresRequiredToFail: number
}

export enum Allegiance {
  Resistance = "Resistance",
  Spies = "Spies",
}