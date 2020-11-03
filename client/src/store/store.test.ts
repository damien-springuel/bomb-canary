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
      player: "",
      players: [],
      leader: "",
      isPlayerTheLeader: false,
    }
  );
});

test(`Store - endReplay`, t => {
  const store = new Store();
  store.startReplay();
  store.showPartyRoom("test");
  store.assignLeader("leader");
  
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Loading);
  t.deepEqual(storeValues.partyCode, "");
  t.deepEqual(storeValues.leader, "");

  store.endReplay()
  
  storeValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.PartyRoom);
  t.deepEqual(storeValues.partyCode, "test");
  t.deepEqual(storeValues.leader, "leader");
});

test(`Store - endReplay twice`, t => {
  const store = new Store();
  store.startReplay();
  store.showPartyRoom("test");
  store.assignLeader("leader")
  
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Loading);
  t.deepEqual(storeValues.partyCode, "");
  t.deepEqual(storeValues.leader, "");

  store.endReplay()
  store.endReplay()

  storeValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.PartyRoom);
  t.deepEqual(storeValues.partyCode, "test");
  t.deepEqual(storeValues.leader, "leader");
});

test(`Store - startReplay twice`, t => {
  const store = new Store();
  
  store.startReplay();
  store.startReplay();
  store.showPartyRoom("test");
  store.assignLeader("leader");
  
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Loading);
  t.deepEqual(storeValues.partyCode, "");

  store.endReplay()
  
  storeValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.PartyRoom);
  t.deepEqual(storeValues.partyCode, "test");
  t.deepEqual(storeValues.leader, "leader");
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

test(`Store - showLobby`, t => {
  const store = new Store();
  store.showLobby();
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Lobby);
});

test(`Store - showPartyRoom`, t => {
  const store = new Store();
  store.showPartyRoom("testCode");
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.PartyRoom);
  t.deepEqual(storeValues.partyCode, "testCode");
});

test(`Store - showGameRoom`, t => {
  const store = new Store();
  store.showGameRoom();
  const storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.pageToShow, Page.Game);
});

test(`Store - definePlayer`, t => {
  const store = new Store();
  store.definePlayer("testName");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.player, "testName");
});

test(`Store - joinPlayer`, t => {
  const store = new Store();
  store.joinPlayer("testName1");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.players, ["testName1"]);


  store.joinPlayer("testName2");
  storeValues = get(store);
  t.deepEqual(storeValues.players, ["testName1", "testName2"]);
});

test(`Store - assignLeader`, t => {
  const store = new Store();
  store.assignLeader("testName1");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.leader, "testName1");
});

test(`Store - isLeader`, t => {
  const store = new Store();
  store.definePlayer("testName");
  store.assignLeader("anotherLeader");
  let storeValues: StoreValues = get(store);
  t.false(storeValues.isPlayerTheLeader);
  
  store.assignLeader("testName");
  storeValues = get(store);
  t.true(storeValues.isPlayerTheLeader);
});