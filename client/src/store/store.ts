import {get, type Readable, type Writable} from "svelte/store";
import {writable} from "svelte/store";
import { 
  Dialog,
  GamePhase,
  Page, 
  type MissionRequirement, 
  type MissionResult, 
  type TeamVotes 
} from "../types/types";

export interface StoreValues {
  pageToShow: Page
  partyCode: string
  player: string
  players: string[]
  missionRequirements: MissionRequirement[]
  currentMission: number,
  currentGamePhase: GamePhase,
  leader: string
  currentTeam: Set<string>,
  peopleThatVotedOnTeam: Set<string>,
  playerVote: boolean | null
  currentTeamVoteNb: number,
  teamVoteResults: TeamVotes[],
  peopleThatWorkedOnMission: Set<string>,
  playerMissionSuccess: boolean | null
  missionResults: MissionResult[]
  dialogShown: Dialog,
  revealedSpies: Set<string>,
  missionDetailsShown: number
}

function defaultValues(): StoreValues {
  return {
    pageToShow: Page.Loading,
    partyCode: "",
    player: "",
    players: [],
    missionRequirements: [],
    currentMission: 0,
    currentGamePhase: GamePhase.TeamSelection,
    leader: "",
    currentTeam: new Set<string>(),
    peopleThatVotedOnTeam: new Set<string>(),
    playerVote: null,
    currentTeamVoteNb: 1,
    teamVoteResults: [{votes: []}, {votes: []}, {votes: []}, {votes: []}, {votes: []}],
    peopleThatWorkedOnMission: new Set<string>(),
    playerMissionSuccess: null,
    missionResults: [],
    dialogShown: null,
    revealedSpies: new Set<string>(),
    missionDetailsShown: 0,
  }
}

export class Store implements Readable<StoreValues> {
  protected replayingEvent: boolean = false;
  protected replayedValues: StoreValues = defaultValues();
  protected readonly writable: Writable<StoreValues> = writable(defaultValues());

  subscribe(run: (value: StoreValues) => void, invalidate?: (value?: StoreValues) => void): () => void {
    return this.writable.subscribe(run, invalidate);
  }

  protected update(updater: (value: StoreValues) => StoreValues) {
    if (this.replayingEvent) {
      this.replayedValues = this.updateComputed(updater(this.replayedValues));
    } 
    else {
      this.writable.update(v => this.updateComputed(updater(v)));
    }
  }

  protected updateNoReplay(updater: (value: StoreValues) => StoreValues) {
    if (!this.replayingEvent) {
      this.writable.update(v => this.updateComputed(updater(v)));
    }
  }

  protected updateComputed(value: StoreValues): StoreValues {
    value.currentMission = value.missionResults.length;
    value.currentTeamVoteNb = value.teamVoteResults[value.currentMission].votes.length + 1;
    
    return value;
  }

  startReplay() {
    if (!this.replayingEvent) {
      this.replayedValues = {...get(this.writable)};
      this.replayingEvent = true;
    }
  }

  endReplay() {
    if (this.replayingEvent) {
      this.writable.set({...this.replayedValues});
      this.replayedValues = null;
      this.replayingEvent = false;
    }
  }

  reset() {
    this.writable.set(defaultValues());
  }

  readonly showLobby = showLobby;
  readonly showPartyRoom = showPartyRoom;
  readonly showGameRoom = showGameRoom;
  readonly joinPlayer = joinPlayer;
  readonly definePlayer = definePlayer;
  readonly setMissionRequirements = setMissionRequirements;
  readonly startTeamSelection = startTeamSelection;
  readonly assignLeader = assignLeader;
  readonly selectPlayer = selectPlayer;
  readonly deselectPlayer = deselectPlayer;
  readonly startTeamVote = startTeamVote;
  readonly makePlayerVote = makePlayerVote;
  readonly saveTeamVoteResult = saveTeamVoteResult;
  readonly startMission = startMission;
  readonly makePlayerWorkOnMission = makePlayerWorkOnMission;
  readonly saveMissionResult = saveMissionResult;
  readonly showIdentity = showIdentity;
  readonly showMissionDetails = showMissionDetails;
  readonly closeDialog = closeDialog;
  readonly rememberSpies = rememberSpies;
}

function showLobby(this: Store) {
  this.update(v => {
    v.pageToShow = Page.Lobby
    return v;
  });
}

function showPartyRoom(this: Store, code: string) {
  this.update(v => {
    v.pageToShow = Page.PartyRoom
    v.partyCode = code;
    return v;
  });
}

function showGameRoom(this: Store) {
  this.update(v => {
    v.pageToShow = Page.Game
    return v;
  });
}

function definePlayer(this: Store, name: string) {
  this.update(v => {
    v.player = name;
    return v;
  });
}

function joinPlayer(this: Store, name: string) {
  this.update(v => {
    v.players.push(name);
    return v;
  });
}

function setMissionRequirements(this: Store, requirements: MissionRequirement[]) {
  this.update(v => {
    v.missionRequirements = requirements.slice();
    return v;
  });
}

function startTeamSelection(this: Store): void {
  this.update(v => {
    v.currentGamePhase = GamePhase.TeamSelection;
    v.currentTeam.clear();
    v.peopleThatVotedOnTeam.clear();
    v.playerVote = null;
    v.peopleThatWorkedOnMission.clear();
    v.playerMissionSuccess = null;
    return v;
  });
}

function assignLeader(this: Store, leader: string) {
  this.update(v => {
    v.leader = leader;
    return v;
  });
}

function selectPlayer(this: Store, player: string) {
  this.update(v => {
    v.currentTeam.add(player);
    return v;
  });
}

function deselectPlayer(this: Store, player: string) {
  this.update(v => {
    v.currentTeam.delete(player);
    return v;
  });
}

function startTeamVote(this: Store): void {
  this.update(v => {
    v.currentGamePhase = GamePhase.TeamVote;
    return v;
  });
}

function makePlayerVote(this: Store, player: string, approval: boolean | null): void {
  this.update(v => {
    v.peopleThatVotedOnTeam.add(player);
    if (player === v.player) {
      v.playerVote = approval;
    }
    return v;
  });
}

function saveTeamVoteResult(this: Store, approved: boolean, playerVotes: Map<string, boolean>): void {
  this.update(v => {
    v.teamVoteResults[v.currentMission].votes.push({
      team: new Set<string>(v.currentTeam),
      approved: approved, 
      playerVotes: playerVotes
    });
    return v;
  });
}

function startMission(this: Store): void {
  this.update(v => {
    v.currentGamePhase = GamePhase.Mission;
    return v;
  });
}

function makePlayerWorkOnMission(this: Store, player: string, success: boolean | null): void {
  this.update(v => {
    v.peopleThatWorkedOnMission.add(player);
    if (player === v.player) {
      v.playerMissionSuccess = success;
    }
    return v;
  });
}

function saveMissionResult(this: Store, success: boolean, nbFails: number): void {
  this.update(v => {
    v.missionResults.push({success, nbFails});
    return v;
  });
}

function showIdentity(this: Store) {
  this.updateNoReplay(v => {
    v.dialogShown = Dialog.Identity;
    return v
  })
}

function showMissionDetails(this: Store, mission:number) {
  this.updateNoReplay(v => {
    v.dialogShown = Dialog.MissionDetails;
    v.missionDetailsShown = mission;
    return v
  })
}

function closeDialog(this: Store) {
  this.update(v => {
    v.dialogShown = null;
    return v
  })
}

function rememberSpies(this: Store, spies: Set<string>) {
  this.update(v => {
    v.revealedSpies = spies;
    return v
  })
}
