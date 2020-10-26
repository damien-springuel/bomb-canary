import test from "ava";
import { Page, Store, StoreValues } from "./store";
import {get} from "svelte/store";

test(`Store - default values`, t => {
  const store = new Store();
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues, 
    {
      pageToShow: Page.Loading,
      partyCode: "",
      players: [],
    }
  );
});

test(`Store - endReplay`, t => {
  const store = new Store();
  store.showPartyRoom("test");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Loading);
  t.deepEqual(storeValues.partyCode, "");

  store.endReplay()
  storeValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.PartyRoom);
  t.deepEqual(storeValues.partyCode, "test");
});

test(`Store - endReplay twice`, t => {
  const store = new Store();
  store.showPartyRoom("test");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Loading);
  t.deepEqual(storeValues.partyCode, "");

  store.endReplay()
  store.endReplay()
  storeValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.PartyRoom);
  t.deepEqual(storeValues.partyCode, "test");
});

test(`Store - reset`, t => {
  const store = new Store();
  store.showPartyRoom("test");
  store.joinPlayer("name1");
  store.reset();
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Loading);
  t.deepEqual(storeValues.players, []);
});

function getReplayEndedStore(): Store {
  const store = new Store();
  store.endReplay();
  return store;
}

test(`Store - showLobby`, t => {
  const store = getReplayEndedStore();
  store.showLobby();
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Lobby);
});

test(`Store - showPartyRoom`, t => {
  const store = getReplayEndedStore();
  store.showPartyRoom("testCode");
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.PartyRoom);
  t.deepEqual(storeValues.partyCode, "testCode");
});

test(`Store - showGameRoom`, t => {
  const store = getReplayEndedStore();
  store.showGameRoom();
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Game);
});

test(`Store - joinPlayer`, t => {
  const store = getReplayEndedStore();
  store.joinPlayer("testName1");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.players, ["testName1"]);


  store.joinPlayer("testName2");
  storeValues = get(store);
  t.deepEqual(storeValues.players, ["testName1", "testName2"]);
});