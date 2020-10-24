import type {Writable, Readable} from "svelte/store";
import {writable} from "svelte/store";

export enum Page {
  Loading = "loading",
  Lobby = "lobby",
  PartyRoom = "partyRoom",
}

export interface StoreValues {
  pageToShow: Page
  partyCode: string
}

function defaultValues(): StoreValues {
  return {
    pageToShow: Page.Loading,
    partyCode: "",
  }
}

export class Store implements Readable<StoreValues> {

  protected replayingEvent: boolean = true;
  protected replayedValues: StoreValues = defaultValues();
  protected readonly writable: Writable<StoreValues> = writable(defaultValues());

  subscribe(run: (value: StoreValues) => void, invalidate?: (value?: StoreValues) => void): () => void {
    return this.writable.subscribe(run, invalidate);
  }

  protected update(updater: (value: StoreValues) => StoreValues) {
    if (this.replayingEvent) {
      this.replayedValues = updater(this.replayedValues);
    } 
    else {
      this.writable.update(updater);
    }
  }

  endReplay() {
    if (this.replayingEvent) {
      this.writable.set({...this.replayedValues});
      this.replayedValues = null;
      this.replayingEvent = false;
    }
  }

  readonly showLobby = showLobby;
  readonly showPartyRoom = showPartyRoom;
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