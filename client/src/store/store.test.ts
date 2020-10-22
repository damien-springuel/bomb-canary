import test from "ava";
import { Page, Store, StoreValues } from "./store";
import {get} from "svelte/store";

test(`Store - default values`, t => {
  const store = new Store();
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues, 
    {
      pageToShow: Page.Loading,
    }
  );
});

test(`Store - showLobby`, t => {
  const store = new Store();
  store.showLobby();
  const storeValues: StoreValues = get(store);
  t.is(storeValues.pageToShow, Page.Lobby);
});