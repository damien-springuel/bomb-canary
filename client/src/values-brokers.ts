import type { GameValues } from "./components/Game-service";
import type { IdentityValues } from "./components/Identity-service";
import type { MissionConductingValues } from "./components/MissionConducting-service";
import type { MissionDetailsValues } from "./components/MissionDetails-service";
import type { MissionTrackerValues } from "./components/MissionTracker-service";
import type { PageValues } from "./components/Page-service";
import type { PartyRoomValues } from "./components/PartyRoom-service";
import type { TeamSelectionValues } from "./components/TeamSelection-service";
import type { TeamVoteValues } from "./components/TeamVote-service";
import type { StoreValues } from "./store/store";
import type { Dialog, GamePhase, MissionRequirement, MissionResult, Page, TeamVotes } from "./types/types";

export class IdentityValuesBroker implements IdentityValues {
  constructor(private readonly storeValues: StoreValues){}

  get player(): string {
    return this.storeValues.player;
  }
  get revealedSpies(): Set<string>{
    return this.storeValues.revealedSpies;
  }
}

export class MissionTrackerValuesBroker implements MissionTrackerValues {
  constructor(private readonly storeValues: StoreValues) {}
  
  get missionRequirements(): MissionRequirement[] {
    return this.storeValues.missionRequirements;
  }

  get missionResults(): MissionResult[] {
    return this.storeValues.missionResults;
  }
}

export class MissionDetailsValuesBroker implements MissionDetailsValues {
  constructor(private readonly storeValues: StoreValues) {}
  get mission(): number{
    return this.storeValues.missionDetailsShown;
  }

  get teamVotes(): TeamVotes {
    return this.storeValues.teamVoteResults[this.mission];
  }
}


export class TeamSelectionValuesBroker implements TeamSelectionValues {
  constructor(private readonly storeValues: StoreValues) {}
  
  get currentTeam(): Set<string> {
    return this.storeValues.currentTeam;
  }
  
  get player(): string {
    return this.storeValues.player;
  }
  
  get leader(): string {
    return this.storeValues.leader;
  }
  
  get currentTeamVoteNb(): number {
    return this.storeValues.currentTeamVoteNb;
  }
  
  get players(): string[] {
    return this.storeValues.players;
  }

  get nbPeopleRequiredOnMission(): number {
    return this.storeValues.missionRequirements[this.storeValues.currentMission].nbPeopleOnMission;
  }
}

export class TeamVoteValuesBroker implements TeamVoteValues {
  constructor(private readonly storeValues: StoreValues){}
  
  get player(): string {
    return this.storeValues.player;
  }
  
  get players(): string[] {
    return this.storeValues.players;
  }
  
  get currentTeam(): Set<string> {
    return this.storeValues.currentTeam;
  }
  
  get peopleThatVotedOnTeam(): Set<string> {
    return this.storeValues.peopleThatVotedOnTeam;
  }
  
  get playerVote(): boolean {
    return this.storeValues.playerVote;
  }
}

export class MissionConductingValuesBroker implements MissionConductingValues {
  constructor(private readonly storeValues: StoreValues){}
  
  get player(): string {
    return this.storeValues.player;
  }
  
  get currentTeam(): Set<string> {
    return this.storeValues.currentTeam;
  }
  
  get peopleThatWorkedOnMission(): Set<string> {
    return this.storeValues.peopleThatWorkedOnMission;
  }
  
  get playerMissionSuccess(): boolean {
    return this.storeValues.playerMissionSuccess;
  }
}

export class GameValuesBroker implements GameValues {
  constructor(private readonly storeValues: StoreValues){}
    
  get currentGamePhase(): GamePhase {
    return this.storeValues.currentGamePhase;
  }
  
  get dialogShown(): Dialog {
    return this.storeValues.dialogShown;
  };
  
  get identityValues(): IdentityValues {
    return new IdentityValuesBroker(this.storeValues);
  }
  get missionTrackerValues(): MissionTrackerValues {
    return new MissionTrackerValuesBroker(this.storeValues);
  }

  get teamSelectionValues(): TeamSelectionValues {
    return new TeamSelectionValuesBroker(this.storeValues);
  }

  get teamVoteValues(): TeamVoteValues {
    return new TeamVoteValuesBroker(this.storeValues);
  }

  get missionConductingValues(): MissionConductingValues {
    return new MissionConductingValuesBroker(this.storeValues);
  }

  get missionDetailsValues(): MissionDetailsValues {
    return new MissionDetailsValuesBroker(this.storeValues);
  }
}

export class PartyRoomValuesBroker implements PartyRoomValues {
  constructor(private readonly storeValues: StoreValues){}
  
  get partyCode(): string {
    return this.storeValues.partyCode;
  }
  
  get players(): string[] {
    return this.storeValues.players;
  }
}

export class PageValuesBroker implements PageValues {
  constructor(private readonly storeValues: StoreValues){}
  
  get pageToShow(): Page {
    return this.storeValues.pageToShow;
  };
  
  get gameValues(): GameValues {
    return new GameValuesBroker(this.storeValues);
  }
  
  get partyRoomValues(): PartyRoomValues {
    return new PartyRoomValuesBroker(this.storeValues);
  }
}

export class AppValuesBroker {
  constructor(private readonly storeValues: StoreValues) {}

  get pageValues(): PageValues {
    return new PageValuesBroker(this.storeValues);
  }
}