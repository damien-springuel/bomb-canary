import test from "ava";
import { GamePhase, Page, Store, StoreValues } from "./store";
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
      missionRequirements: [],
      currentMission: 1,
      currentGamePhase: GamePhase.TeamSelection,
      leader: "",
      isPlayerTheLeader: false,
      currentTeam: new Set<string>(),
      isPlayerInTeam: undefined,
      isPlayerSelectableForTeam: undefined,
      canConfirmTeam: false,
      peopleThatVotedOnTeam: new Set<string>(),
      playerVote: null,
      hasGivenPlayerVoted: undefined,
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

test(`Store - setMissionRequirements`, t => {
  const store = new Store();
  store.setMissionRequirements([{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.missionRequirements, [{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  t.deepEqual(storeValues.currentMission, 1);
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

test(`Store - selectPlayer`, t => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.currentTeam, new Set<string>(["p1", "p2"]));
});

test(`Store - deselectPlayer`, t => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  store.deselectPlayer("p1");
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.currentTeam, new Set<string>(["p2"]));
});

test(`Store - isPlayerInTeam`, t => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  store.deselectPlayer("p1");
  let storeValues: StoreValues = get(store);
  t.false(storeValues.isPlayerInTeam("p1"));
  t.true(storeValues.isPlayerInTeam("p2"));
  t.false(storeValues.isPlayerInTeam("p3"));
});

test(`Store - isPlayerSelectableForTeam`, t => {
  const store = new Store();
  store.setMissionRequirements([{nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1}]);
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  let storeValues: StoreValues = get(store);
  t.false(storeValues.isPlayerSelectableForTeam("p3"));
  t.true(storeValues.isPlayerSelectableForTeam("p1"));
  
  store.deselectPlayer("p1");
  storeValues = get(store);
  t.true(storeValues.isPlayerSelectableForTeam("p3"));
});

test(`Store - canConfirmTeam`, t => {
  const store = new Store();
  store.setMissionRequirements([{nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1}]);
  store.selectPlayer("p1");
  let storeValues: StoreValues = get(store);
  t.false(storeValues.canConfirmTeam);
  
  store.selectPlayer("p2");
  storeValues = get(store);
  t.true(storeValues.canConfirmTeam);
});

test(`Store - startTeamVote`, t => {
  const store = new Store();
  store.startTeamVote();
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.currentGamePhase, GamePhase.TeamVote);
});

test(`Store - makePlayerVote - not the player`, t => {
  const store = new Store();
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.peopleThatVotedOnTeam, new Set<string>(["testName"]));
  t.deepEqual(storeValues.playerVote, null);
});

test(`Store - makePlayerVote - the player`, t => {
  const store = new Store();
  store.definePlayer("testName");
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  t.deepEqual(storeValues.peopleThatVotedOnTeam, new Set<string>(["testName"]));
  t.deepEqual(storeValues.playerVote, true);
});

test(`Store - hasPlayerVoted`, t => {
  const store = new Store();
  store.makePlayerVote("p1", false);
  let storeValues: StoreValues = get(store);
  t.true(storeValues.hasGivenPlayerVoted("p1"));
  t.false(storeValues.hasGivenPlayerVoted("p2"));
});