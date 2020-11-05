import {Writable, Readable, get} from "svelte/store";
import {writable} from "svelte/store";
import type { MissionRequirement } from "../messages/events";

export enum Page {
  Loading = "loading",
  Lobby = "lobby",
  PartyRoom = "partyRoom",
  Game = "game",
}

export interface StoreValues {
  pageToShow: Page
  partyCode: string
  player: string
  players: string[]
  missionRequirements: MissionRequirement[]
  currentMission: number,
  leader: string
  isPlayerTheLeader: boolean
  currentTeam: Set<string>,
  isPlayerInTeam: (player: string) => boolean,
  isPlayerSelectableForTeam: (player: string) => boolean,
  canConfirmTeam: boolean,
}

function defaultValues(): StoreValues {
  return {
    pageToShow: Page.Loading,
    partyCode: "",
    player: "",
    players: [],
    missionRequirements: [],
    currentMission: 1,
    leader: "",
    isPlayerTheLeader: false,
    currentTeam: new Set<string>(),
    isPlayerInTeam: undefined,
    isPlayerSelectableForTeam: undefined,
    canConfirmTeam: false,
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
    value.isPlayerInTeam = player => value.currentTeam.has(player);
    
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