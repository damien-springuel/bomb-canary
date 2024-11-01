import { expect, test } from "vitest";
import { Store, type StoreValues } from "./store";
import {get} from "svelte/store";
import { Allegiance, Dialog, GamePhase, Page } from "../types/types";

test(`default values`, () => {
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
      winner: null,
    }
  );
});

test(`endReplay`, () => {
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

test(`endReplay twice`, () => {
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

test(`startReplay twice`, () => {
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

test(`reset`, () => {
  const store = new Store();
  store.showPartyRoom("test");
  store.joinPlayer("name1");
  store.reset();
  let storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Loading);
  expect(storeValues.players).to.deep.equal([]);
});

test(`showLobby`, () => {
  const store = new Store();
  store.showLobby();
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Lobby);
});

test(`showPartyRoom`, () => {
  const store = new Store();
  store.showPartyRoom("testCode");
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.PartyRoom);
  expect(storeValues.partyCode).to.equal("testCode");
});

test(`showGameRoom`, () => {
  const store = new Store();
  store.showGameRoom();
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Game);
});

test(`definePlayer`, () => {
  const store = new Store();
  store.definePlayer("testName");
  let storeValues: StoreValues = get(store);
  expect(storeValues.player).to.equal("testName");
});

test(`joinPlayer`, () => {
  const store = new Store();
  store.joinPlayer("testName1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.players).to.deep.equal(["testName1"]);


  store.joinPlayer("testName2");
  storeValues = get(store);
  expect(storeValues.players).to.deep.equal(["testName1", "testName2"]);
});

test(`setMissionRequirements`, () => {
  const store = new Store();
  store.setMissionRequirements([{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  let storeValues: StoreValues = get(store);
  expect(storeValues.missionRequirements).to.deep.equal([{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  expect(storeValues.currentMission).to.equal(0);
});

test(`startTeamSelection`, () => {
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

test(`assignLeader`, () => {
  const store = new Store();
  store.assignLeader("testName1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.leader).to.equal("testName1");
});

test(`selectPlayer`, () => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeam).to.deep.equal(new Set<string>(["p1", "p2"]));
});

test(`deselectPlayer`, () => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  store.deselectPlayer("p1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeam).to.deep.equal(new Set<string>(["p2"]));
});

test(`startTeamVote`, () => {
  const store = new Store();
  store.startTeamVote();
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentGamePhase).to.equal(GamePhase.TeamVote);
});

test(`makePlayerVote - not the player`, () => {
  const store = new Store();
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatVotedOnTeam).to.deep.equal(new Set<string>(["testName"]));
  expect(storeValues.playerVote).to.be.null;
});

test(`makePlayerVote - the player`, () => {
  const store = new Store();
  store.definePlayer("testName");
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatVotedOnTeam).to.deep.equal(new Set<string>(["testName"]));
  expect(storeValues.playerVote).to.be.true;
});

test(`saveTeamVoteResult`, () => {
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

test(`saveTeamVoteResult - multiple results`, () => {
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

test(`saveTeamVoteResult - multiple results in multiple missions`, () => {
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

test(`currentTeamVoteNb`, () => {
  const store = new Store();
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  store.saveTeamVoteResult(false, new Map<string, boolean>([["p1", false], ["p2", false]]));
  store.saveMissionResult(false, 2);
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", true]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeamVoteNb).to.equal(2);
});

test(`startMission`, () => {
  const store = new Store();
  store.startMission();
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentGamePhase).to.equal(GamePhase.Mission);
});

test(`makePlayerWorkOnMission - the player`, () => {
  const store = new Store();
  store.definePlayer("p1");
  store.makePlayerWorkOnMission("p1", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatWorkedOnMission).to.deep.equal(new Set<string>(["p1"]));
  expect(storeValues.playerMissionSuccess).to.be.true;
});

test(`makePlayerWorkOnMission - not the player`, () => {
  const store = new Store();
  store.definePlayer("p1");
  store.makePlayerWorkOnMission("p2", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatWorkedOnMission).to.deep.equal(new Set<string>(["p2"]));
  expect(storeValues.playerMissionSuccess).to.be.null;
});

test(`saveMissionResult`, () => {
  const store = new Store();
  store.saveMissionResult(true, 2);
  let storeValues: StoreValues = get(store);
  expect(storeValues.missionResults).to.deep.equal([{success: true, nbFails: 2}]);
  expect(storeValues.currentMission).to.equal(1); // second mission, since mission are zero-based
});

test(`close Dialog`, () => {
  const store = new Store();
  store.showIdentity();
  store.closeDialog();
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`showIdentity`, () => {
  const store = new Store();
  store.showIdentity();
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.equal(Dialog.Identity);
});

test(`showIdentity isn't replayed`, () => {
  const store = new Store();
  
  store.startReplay();
  store.showIdentity();
  store.endReplay();

  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`showMissionDetails`, () => {
  const store = new Store();
  store.showMissionDetails(2);
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.equal(Dialog.MissionDetails);
  expect(storeValues.missionDetailsShown).to.equal(2);
});

test(`showMissionDetails isn't replayed`, () => {
  const store = new Store();
  
  store.startReplay();
  store.showMissionDetails(2);
  store.endReplay();

  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`rememberSpies`, () => {
  const store = new Store();
  store.rememberSpies(new Set<string>(["spy 1", "spy 2"]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.revealedSpies).to.deep.equal(new Set<string>(["spy 1", "spy 2"]));
});

test(`showLastMissionResult`, () => {
  const store = new Store();
  store.showLastMissionResult();
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.equal(Dialog.LastMissionResult);
});

test(`showLastMissionResult isn't replayed`, () => {
  const store = new Store();
  
  store.startReplay();
  store.showLastMissionResult();
  store.endReplay();

  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`end game`, () => {
  const store = new Store();
  store.endGame(Allegiance.Resistance, new Set<string>(["spy 1", "spy2"]));

  let storeValues: StoreValues = get(store);
  expect(storeValues.winner).to.equal(Allegiance.Resistance);
  expect(storeValues.revealedSpies).to.deep.equal(new Set<string>(["spy 1", "spy2"]));
  expect(storeValues.currentGamePhase).to.equal(GamePhase.GameEnded);
});