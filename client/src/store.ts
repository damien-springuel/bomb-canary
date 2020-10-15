import type {Writable, Readable} from "svelte/store";
import {writable} from "svelte/store";

interface StoreValues {
  name: string
}

export class Store implements Readable<StoreValues> {

  protected readonly writable: Writable<StoreValues>;
  constructor() {
    this.writable = writable({} as StoreValues);
  }

  subscribe(run: (value: StoreValues) => void, invalidate?: (value?: StoreValues) => void): () => void {
    return this.writable.subscribe(run, invalidate);
  }

  setName = setName;
}

function setName(this: Store, name: string) {
  this.writable.update(v => {
    v.name = name
    return v;
  });
}