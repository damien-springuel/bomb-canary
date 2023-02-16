import { expect, test } from "vitest";
import { Dialog, GamePhase, Page, Store, type StoreValues } from "./store";
import {get} from "svelte/store";

test(`Store - default values`, t => {
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
      isPlayerTheLeader: false,
      currentTeam: new Set<string>(),
      isGivenPlayerInTeam: undefined,
      isPlayerSelectableForTeam: undefined,
      canConfirmTeam: false,
      peopleThatVotedOnTeam: new Set<string>(),
      playerVote: null,
      hasGivenPlayerVoted: undefined,
      currentTeamVoteNb: 1,
      teamVoteResults: [{votes: []}, {votes: []}, {votes: []}, {votes: []}, {votes: []}],
      isPlayerInMission: false,
      peopleThatWorkedOnMission: new Set<string>(),
      playerMissionSuccess: null,
      hasGivenPlayerWorkedOnMission: undefined,
      missionResults: [],
      dialogShown: null,
      revealedSpies: new Set<string>(),
    }
  );
});

test(`Store - endReplay`, t => {
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

test(`Store - endReplay twice`, t => {
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

test(`Store - startReplay twice`, t => {
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

test(`Store - reset`, t => {
  const store = new Store();
  store.showPartyRoom("test");
  store.joinPlayer("name1");
  store.reset();
  let storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Loading);
  expect(storeValues.players).to.deep.equal([]);
});

test(`Store - showLobby`, t => {
  const store = new Store();
  store.showLobby();
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Lobby);
});

test(`Store - showPartyRoom`, t => {
  const store = new Store();
  store.showPartyRoom("testCode");
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.PartyRoom);
  expect(storeValues.partyCode).to.equal("testCode");
});

test(`Store - showGameRoom`, t => {
  const store = new Store();
  store.showGameRoom();
  const storeValues: StoreValues = get(store);
  expect(storeValues.pageToShow).to.equal(Page.Game);
});

test(`Store - definePlayer`, t => {
  const store = new Store();
  store.definePlayer("testName");
  let storeValues: StoreValues = get(store);
  expect(storeValues.player).to.equal("testName");
});

test(`Store - joinPlayer`, t => {
  const store = new Store();
  store.joinPlayer("testName1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.players).to.deep.equal(["testName1"]);


  store.joinPlayer("testName2");
  storeValues = get(store);
  expect(storeValues.players).to.deep.equal(["testName1", "testName2"]);
});

test(`Store - setMissionRequirements`, t => {
  const store = new Store();
  store.setMissionRequirements([{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  let storeValues: StoreValues = get(store);
  expect(storeValues.missionRequirements).to.deep.equal([{nbFailuresRequiredToFail: 3, nbPeopleOnMission: 4}, {nbFailuresRequiredToFail: 2, nbPeopleOnMission:4}]);
  expect(storeValues.currentMission).to.equal(0);
});

test(`Store - startTeamSelection`, t => {
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

test(`Store - assignLeader`, t => {
  const store = new Store();
  store.assignLeader("testName1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.leader).to.equal("testName1");
});

test(`Store - isLeader`, t => {
  const store = new Store();
  store.definePlayer("testName");
  store.assignLeader("anotherLeader");
  let storeValues: StoreValues = get(store);
  expect(storeValues.isPlayerTheLeader).to.be.false;
  
  store.assignLeader("testName");
  storeValues = get(store);
  expect(storeValues.isPlayerTheLeader).to.be.true;
});

test(`Store - selectPlayer`, t => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeam).to.deep.equal(new Set<string>(["p1", "p2"]));
});

test(`Store - deselectPlayer`, t => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  store.deselectPlayer("p1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeam).to.deep.equal(new Set<string>(["p2"]));
});

test(`Store - isGivenPlayerInTeam`, t => {
  const store = new Store();
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  store.deselectPlayer("p1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.isGivenPlayerInTeam("p1")).to.be.false;
  expect(storeValues.isGivenPlayerInTeam("p2")).to.be.true;
  expect(storeValues.isGivenPlayerInTeam("p3")).to.be.false;
});

test(`Store - isPlayerSelectableForTeam`, t => {
  const store = new Store();
  store.setMissionRequirements([{nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1}]);
  store.selectPlayer("p1");
  store.selectPlayer("p2");
  let storeValues: StoreValues = get(store);
  expect(storeValues.isPlayerSelectableForTeam("p3")).to.be.false;
  expect(storeValues.isPlayerSelectableForTeam("p1")).to.be.true;
  
  store.deselectPlayer("p1");
  storeValues = get(store);
  expect(storeValues.isPlayerSelectableForTeam("p3")).to.be.true;
});

test(`Store - canConfirmTeam`, t => {
  const store = new Store();
  store.setMissionRequirements([{nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1}]);
  store.selectPlayer("p1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.canConfirmTeam).to.be.false;
  
  store.selectPlayer("p2");
  storeValues = get(store);
  expect(storeValues.canConfirmTeam).to.be.true;
});

test(`Store - startTeamVote`, t => {
  const store = new Store();
  store.startTeamVote();
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentGamePhase).to.equal(GamePhase.TeamVote);
});

test(`Store - makePlayerVote - not the player`, t => {
  const store = new Store();
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatVotedOnTeam).to.deep.equal(new Set<string>(["testName"]));
  expect(storeValues.playerVote).to.be.null;
});

test(`Store - makePlayerVote - the player`, t => {
  const store = new Store();
  store.definePlayer("testName");
  store.makePlayerVote("testName", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatVotedOnTeam).to.deep.equal(new Set<string>(["testName"]));
  expect(storeValues.playerVote).to.be.true;
});

test(`Store - saveTeamVoteResult`, t => {
  const store = new Store();
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.teamVoteResults).to.deep.equal([
    {votes: [{approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])}]}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}
  ]);
});

test(`Store - saveTeamVoteResult - multiple results`, t => {
  const store = new Store();
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  store.saveTeamVoteResult(false, new Map<string, boolean>([["p1", false], ["p2", false]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.teamVoteResults).to.deep.equal([
    {votes: [
      {approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])},
      {approved: false, playerVotes: new Map<string, boolean>([["p1", false], ["p2", false]])},
    ]}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}
  ]);
});

test(`Store - saveTeamVoteResult - multiple results in multiple missions`, t => {
  const store = new Store();
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  store.saveTeamVoteResult(false, new Map<string, boolean>([["p1", false], ["p2", false]]));
  store.saveMissionResult(false, 2);
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", true]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.teamVoteResults).to.deep.equal([
    {votes: [
      {approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", false]])},
      {approved: false, playerVotes: new Map<string, boolean>([["p1", false], ["p2", false]])},
    ]}, 
    {votes: [
      {approved: true, playerVotes: new Map<string, boolean>([["p1", true], ["p2", true]])}
    ]}, 
    {votes: []}, 
    {votes: []}, 
    {votes: []}
  ]);
});

test(`Store - hasGivenPlayerVoted`, t => {
  const store = new Store();
  store.makePlayerVote("p1", false);
  let storeValues: StoreValues = get(store);
  expect(storeValues.hasGivenPlayerVoted("p1")).to.be.true;
  expect(storeValues.hasGivenPlayerVoted("p2")).to.be.false;
});

test(`Store - currentTeamVoteNb`, t => {
  const store = new Store();
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", false]]));
  store.saveTeamVoteResult(false, new Map<string, boolean>([["p1", false], ["p2", false]]));
  store.saveMissionResult(false, 2);
  store.saveTeamVoteResult(true, new Map<string, boolean>([["p1", true], ["p2", true]]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentTeamVoteNb).to.equal(2);
});

test(`Store - startMission`, t => {
  const store = new Store();
  store.startMission();
  let storeValues: StoreValues = get(store);
  expect(storeValues.currentGamePhase).to.equal(GamePhase.Mission);
});

test(`Store - isPlayerInMission`, t => {
  const store = new Store();
  store.definePlayer("p1");
  store.selectPlayer("p1");
  let storeValues: StoreValues = get(store);
  expect(storeValues.isPlayerInMission).to.be.true;
  
  store.deselectPlayer("p1");
  storeValues = get(store);
  expect(storeValues.isPlayerInMission).to.be.false;
});

test(`Store - makePlayerWorkOnMission - the player`, t => {
  const store = new Store();
  store.definePlayer("p1");
  store.makePlayerWorkOnMission("p1", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatWorkedOnMission).to.deep.equal(new Set<string>(["p1"]));
  expect(storeValues.playerMissionSuccess).to.be.true;
});

test(`Store - makePlayerWorkOnMission - not the player`, t => {
  const store = new Store();
  store.definePlayer("p1");
  store.makePlayerWorkOnMission("p2", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.peopleThatWorkedOnMission).to.deep.equal(new Set<string>(["p2"]));
  expect(storeValues.playerMissionSuccess).to.be.null;
});

test(`Store - hasGivenPlayerWorkedOnMission`, t => {
  const store = new Store();
  store.makePlayerWorkOnMission("p1", true);
  let storeValues: StoreValues = get(store);
  expect(storeValues.hasGivenPlayerWorkedOnMission("p1")).to.be.true;
  expect(storeValues.hasGivenPlayerWorkedOnMission("p2")).to.be.false;
});

test(`Store - saveMissionResult`, t => {
  const store = new Store();
  store.saveMissionResult(true, 2);
  let storeValues: StoreValues = get(store);
  expect(storeValues.missionResults).to.deep.equal([{success: true, nbFails: 2}]);
  expect(storeValues.currentMission).to.equal(1); // second mission, since mission are zero-based
});

test(`Store - close Dialog`, t => {
  const store = new Store();
  store.showIdentity();
  store.closeDialog();
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`Store - showIdentity`, t => {
  const store = new Store();
  store.showIdentity();
  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.equal(Dialog.Identity);
});

test(`Store - showIdentity isn't replayed`, t => {
  const store = new Store();
  
  store.startReplay();
  store.showIdentity();
  store.endReplay();

  let storeValues: StoreValues = get(store);
  expect(storeValues.dialogShown).to.be.null;
});

test(`Store - rememberSpies`, t => {
  const store = new Store();
  store.rememberSpies(new Set<string>(["spy 1", "spy 2"]));
  let storeValues: StoreValues = get(store);
  expect(storeValues.revealedSpies).to.deep.equal(new Set<string>(["spy 1", "spy 2"]));
});