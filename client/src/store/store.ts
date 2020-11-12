import {Writable, Readable, get} from "svelte/store";
import {writable} from "svelte/store";
import type { MissionRequirement } from "../messages/events";

export enum Page {
  Loading = "loading",
  Lobby = "lobby",
  PartyRoom = "partyRoom",
  Game = "game",
}

export enum GamePhase {
  TeamSelection = "teamSelection",
  TeamVote = "teamVote",
  Mission = "mission",
}

export interface StoreValues {
  pageToShow: Page
  partyCode: string
  player: string
  players: string[]
  missionRequirements: MissionRequirement[]
  currentMission: number,
  currentGamePhase: GamePhase,
  leader: string
  isPlayerTheLeader: boolean
  currentTeam: Set<string>,
  isGivenPlayerInTeam: (player: string) => boolean,
  isPlayerSelectableForTeam: (player: string) => boolean,
  canConfirmTeam: boolean,
  peopleThatVotedOnTeam: Set<string>,
  playerVote: boolean | null
  hasGivenPlayerVoted: (player: string) => boolean,
  isPlayerInMission: boolean;
  peopleThatWorkedOnMission: Set<string>,
  playerMissionSuccess: boolean | null
  hasGivenPlayerWorkedOnMission: (player: string) => boolean,
}

function defaultValues(): StoreValues {
  return {
    pageToShow: Page.Loading,
    partyCode: "",
    player: "",
    players: [],
    missionRequirements: [],
    currentMission: 1,
    currentGamePhase: GamePhase.TeamSelection,
    leader: "",
    isPlayerTheLeader: false,
    currentTeam: new Set<string>(),
    isGivenPlayerInTeam: undefined,
    isPlayerSelectableForTeam: undefined,
    canConfirmTeam: false,
    peopleThatVotedOnTeam: new Set<string>(),
    playerVote: null,
    hasGivenPlayerVoted: undefined,
    isPlayerInMission: false,
    peopleThatWorkedOnMission: new Set<string>(),
    playerMissionSuccess: null,
    hasGivenPlayerWorkedOnMission: undefined,
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

  protected updateComputed(value: StoreValues): StoreValues {
    value.isPlayerTheLeader = !!value.player && !!value.leader && (value.leader === value.player);
    value.isGivenPlayerInTeam = player => value.currentTeam.has(player);
    
    const currentMissionRequirement = value.missionRequirements[value.currentMission-1];
    value.isPlayerSelectableForTeam = player => {
      if (value.currentTeam.has(player)) {
        return true;
      }
      if (value.currentTeam.size < currentMissionRequirement?.nbPeopleOnMission) {
        return true;
      }
      return false;
    }
    value.canConfirmTeam = value.currentTeam.size === currentMissionRequirement?.nbPeopleOnMission
    value.hasGivenPlayerVoted = player => value.peopleThatVotedOnTeam.has(player);
    value.isPlayerInMission = value.currentTeam.has(value.player);
    value.hasGivenPlayerWorkedOnMission = player => value.peopleThatWorkedOnMission.has(player);
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
  readonly assignLeader = assignLeader;
  readonly selectPlayer = selectPlayer;
  readonly deselectPlayer = deselectPlayer;
  readonly startTeamVote = startTeamVote;
  readonly makePlayerVote = makePlayerVote;
  readonly startMission = startMission;
  readonly makePlayerWorkOnMission = makePlayerWorkOnMission;
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
    v.currentMission = 1;
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