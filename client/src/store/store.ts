import {Writable, Readable, get} from "svelte/store";
import {writable} from "svelte/store";

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
  leader: string
  isPlayerTheLeader: boolean
}

function defaultValues(): StoreValues {
  return {
    pageToShow: Page.Loading,
    partyCode: "",
    player: "",
    players: [],
    leader: "",
    isPlayerTheLeader: false,
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
  readonly assignLeader = assignLeader;
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

function assignLeader(this: Store, leader: string) {
  this.update(v => {
    v.leader = leader;
    return v;
  });
}