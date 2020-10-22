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

export class Store implements Readable<StoreValues> {

  protected readonly writable: Writable<StoreValues>;
  constructor() {
    this.writable = writable(
      {
        pageToShow: Page.Loading,
        partyCode: "",
      },
    );
  }

  subscribe(run: (value: StoreValues) => void, invalidate?: (value?: StoreValues) => void): () => void {
    return this.writable.subscribe(run, invalidate);
  }

  showLobby = showLobby;
  showPartyRoom = showPartyRoom;
}

function showLobby(this: Store) {
  this.writable.update(v => {
    v.pageToShow = Page.Lobby
    return v;
  });
}

function showPartyRoom(this: Store, code: string) {
  this.writable.update(v => {
    v.pageToShow = Page.PartyRoom
    v.partyCode = code;
    return v;
  });
}