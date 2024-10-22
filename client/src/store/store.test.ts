import { expect, test } from "vitest";
import { Store, type StoreValues } from "./store";
import {get} from "svelte/store";
import { Dialog, GamePhase, Page } from "../types/types";

test(`Store - default values`, () => {
  const store = new Store();
  const storeValues: StoreValues = get(store);
  expect(storeValues).to.deep.equal( 
    {
      pageToShow: Page.Loading,
      partyCode: "",
      player: "",
      players: [],
      missionRequirements: [],
      currentMission: 0,
      currentGamePhase: GamePhase.TeamSelection,
      leader: "",
      currentTeam: new Set<string>(),
      peopleThatVotedOnTeam: new Set<string>(),
      playerVote: null,
      currentTeamVoteNb: 1,
      teamVoteResults: [{votes: []}, {votes: []}, {votes: []}, {votes: []}, {votes: []}],
      peopleThatWorkedOnMission: new Set<string>(),
      playerMissionSuccess: null,
      missionResults: [],
      dialogShown: null,
      revealedSpies: new Set<string>(),
      missionDetailsShown: 0,
    }
  );
});

test(`Store - endReplay`, () => {
  const store = new Store();
  store.startReplay();
  store.showPartyRoom("test");
  store.assignLeader("leader");
  
  let storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Loading);
  expect(storeValues.partyCode).to.equal("");
  expect(storeValues.leader).to.equal("");

  store.endReplay()
  
  storeValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.PartyRoom);
  expect(storeValues.partyCode).to.equal("test");
  expect(storeValues.leader).to.equal("leader");
});

test(`Store - endReplay twice`, () => {
  const store = new Store();
  store.startReplay();
  store.showPartyRoom("test");
  store.assignLeader("leader")
  
  let storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Loading);
  expect(storeValues.partyCode).to.equal("");
  expect(storeValues.leader).to.equal("");

  store.endReplay()
  store.endReplay()

  storeValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.PartyRoom);
  expect(storeValues.partyCode).to.equal("test");
  expect(storeValues.leader).to.equal("leader");
});

test(`Store - startReplay twice`, () => {
  const store = new Store();
  
  store.startReplay();
  store.startReplay();
  store.showPartyRoom("test");
  store.assignLeader("leader");
  
  let storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Loading);
  expect(storeValues.partyCode).to.equal("");

  store.endReplay()
  
  storeValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.PartyRoom);
  expect(storeValues.partyCode).to.equal("test");
  expect(storeValues.leader).to.equal("leader");
});

test(`Store - reset`, () => {
  const store = new Store();
  store.showPartyRoom("test");
  store.joinPlayer("name1");
  store.reset();
  let storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Loading);
  expect(storeValues.players).to.deep.equal([]);
});

test(`Store - showLobby`, () => {
  const store = new Store();
  store.showLobby();
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Lobby);
});

test(`Store - showPartyRoom`, () => {
  const store = new Store();
  store.showPartyRoom("testCode");
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.PartyRoom);
  expect(storeValues.partyCode).to.equal("testCode");
});

test(`Store - showGameRoom`, () => {
  const store = new Store();
  store.showGameRoom();
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Game);
});

test(`Store - definePlayer`, () => {
  const store = new Store();
  store.definePlayer("testName");
  let storeValues: StoreValues = get(store);
  expect(storeValues.player).to.equal("testName");
});

test(`Store - joinPlayer`, () => {
  const store = new Store();
  store.joinPlayer("testName1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.players).to.deep.equal(["testName1"]);


  store.joinPlayer("testName2");
  storeValues = get(store);
  expect(storeValues.players).to.deep.equal(["testName1", "testName2"]);
});

test(`Store - setMissionRequirements`, () => {
  const store = new Store();
  store.setMissionRequirements([{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  let storeValues: StoreValues = get(store);
  expect(storeValues.missionRequirements).to.deep.equal([{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  expect(storeValues.currentMission).to.equal(0);
});

test(`Store - startTeamSelection`, () => {
  const store = new Store();
  store.definePlayer("testName");
  store.selectPlayer("testName");
  store.makePlayerVote("testName", true);
  store.makePlayerWorkOnMission("testName", true);
  store.startTeamSelection();
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentGamePhase).to.deep.equal(GamePhase.TeamSelection);
  expect(storeValues.peopleThatVotedOnTeam).to.deep.equal(new Set<string>());
  expect(storeValues.playerVote).to.deep.equal(null);
  expect(storeValues.peopleThatWorkedOnMission).to.deep.equal(new Set<string>());
  expect(storeValues.playerMissionSuccess).to.deep.equal(null);
});

test(`Store - assignLeader`, () => {
  const store = new Store();
  store.assignLeader("testName1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.leader).to.equal("testName1");
});

test(`Store - selectPlayer`, () => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeam).to.deep.equal(new Set<string>(["p1", "p2"]));
});

test(`Store - deselectPlayer`, () => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  store.deselectPlayer("p1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeam).to.deep.equal(new Set<string>(["p2"]));
});

test(`Store - startTeamVote`, () => {
  const store = new Store();
  store.startTeamVote();
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentGamePhase).to.equal(GamePhase.TeamVote);
});

test(`Store - makePlayerVote - not the player`, () => {
  const store = new Store();
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatVotedOnTeam).to.deep.equal(new Set<string>(["testName"]));
  expect(storeValues.playerVote).to.be.null;
});

test(`Store - makePlayerVote - the player`, () => {
  const store = new Store();
  store.definePlayer("testName");
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatVotedOnTeam).to.deep.equal(new Set<string>(["testName"]));
  expect(storeValues.playerVote).to.be.true;
});

test(`Store - saveTeamVoteResult`, () => {
  const store = new Store();
  store.selectPlayer("p1");
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.teamVoteResults).to.deep.equal([
    {votes: [{team: new Set<string>(["p1"]), approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])}]}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}
  ]);
});

test(`Store - saveTeamVoteResult - multiple results`, () => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p3");
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  store.saveTeamVoteResult(false, new Map<string, boolean>([["p1", false], ["p2", false]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.teamVoteResults).to.deep.equal([
    {votes: [
      {team: new Set<string>(["p1", "p3"]), approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])},
      {team: new Set<string>(["p1", "p3"]), approved: false, playerVotes: new Map<string, boolean>([["p1", false], ["p2", false]])},
    ]}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}
  ]);
});

test(`Store - saveTeamVoteResult - multiple results in multiple missions`, () => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p3");
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  store.saveTeamVoteResult(false, new Map<string, boolean>([["p1", false], ["p2", false]]));
  store.saveMissionResult(false, 2);
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", true]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.teamVoteResults).to.deep.equal([
    {votes: [
      {team: new Set<string>(["p1", "p3"]), approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])},
      {team: new Set<string>(["p1", "p3"]), approved: false, playerVotes: new Map<string, boolean>([["p1", false], ["p2", false]])},
    ]}, 
    {votes: [
      {team: new Set<string>(["p1", "p3"]), approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", true]])}
    ]}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}
  ]);
});

test(`Store - currentTeamVoteNb`, () => {
  const store = new Store();
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  store.saveTeamVoteResult(false, new Map<string, boolean>([["p1", false], ["p2", false]]));
  store.saveMissionResult(false, 2);
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", true]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeamVoteNb).to.equal(2);
});

test(`Store - startMission`, () => {
  const store = new Store();
  store.startMission();
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentGamePhase).to.equal(GamePhase.Mission);
});

test(`Store - makePlayerWorkOnMission - the player`, () => {
  const store = new Store();
  store.definePlayer("p1");
  store.makePlayerWorkOnMission("p1", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatWorkedOnMission).to.deep.equal(new Set<string>(["p1"]));
  expect(storeValues.playerMissionSuccess).to.be.true;
});

test(`Store - makePlayerWorkOnMission - not the player`, () => {
  const store = new Store();
  store.definePlayer("p1");
  store.makePlayerWorkOnMission("p2", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatWorkedOnMission).to.deep.equal(new Set<string>(["p2"]));
  expect(storeValues.playerMissionSuccess).to.be.null;
});

test(`Store - saveMissionResult`, () => {
  const store = new Store();
  store.saveMissionResult(true, 2);
  let storeValues: StoreValues = get(store);
  expect(storeValues.missionResults).to.deep.equal([{success: true, nbFails: 2}]);
  expect(storeValues.currentMission).to.equal(1); // second mission, since mission are zero-based
});

test(`Store - close Dialog`, () => {
  const store = new Store();
  store.showIdentity();
  store.closeDialog();
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`Store - showIdentity`, () => {
  const store = new Store();
  store.showIdentity();
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.equal(Dialog.Identity);
});

test(`Store - showIdentity isn't replayed`, () => {
  const store = new Store();
  
  store.startReplay();
  store.showIdentity();
  store.endReplay();

  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`Store - showMissionDetails`, () => {
  const store = new Store();
  store.showMissionDetails(2);
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.equal(Dialog.MissionDetails);
  expect(storeValues.missionDetailsShown).to.equal(2);
});

test(`Store - showMissionDetails isn't replayed`, () => {
  const store = new Store();
  
  store.startReplay();
  store.showMissionDetails(2);
  store.endReplay();

  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`Store - rememberSpies`, () => {
  const store = new Store();
  store.rememberSpies(new Set<string>(["spy 1", "spy 2"]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.revealedSpies).to.deep.equal(new Set<string>(["spy 1", "spy 2"]));
});