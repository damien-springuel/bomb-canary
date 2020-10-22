import type {Writable, Readable} from "svelte/store";
import {writable} from "svelte/store";

export enum Page {
  Loading = "loading",
  Lobby = "lobby",
}

export interface StoreValues {
  pageToShow: Page
}

export class Store implements Readable<StoreValues> {

  protected readonly writable: Writable<StoreValues>;
  constructor() {
    this.writable = writable(
      {
        pageToShow: Page.Loading,
      },
    );
  }

  subscribe(run: (value: StoreValues) => void, invalidate?: (value?: StoreValues) => void): () => void {
    return this.writable.subscribe(run, invalidate);
  }

  showLobby = showLobby;
}

function showLobby(this: Store) {
  this.writable.update(v => {
    v.pageToShow = Page.Lobby
    return v;
  });
}