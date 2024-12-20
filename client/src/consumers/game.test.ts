import { expect, test } from "vitest";
import { 
    AllPlayerVotedOnTeam, 
    GameEnded, 
    GameStarted, 
    LeaderConfirmedTeam, 
    LeaderDeselectedMember, 
    LeaderSelectedMember, 
    LeaderStartedToSelectMembers, 
    MissionCompleted, 
    MissionStarted, 
    PlayerVotedOnTeam, 
    PlayerWorkedOnMission 
} from "../messages/events";
import { Allegiance, type MissionRequirement } from "../types/types";
import { GameConsumer, type GameStore } from "./game";

test(`GameStarted`, () => {
  let receivedReq: MissionRequirement[]
  const gameConsumer = new GameConsumer({setMissionRequirements: req => {receivedReq = req}} as GameStore);
  gameConsumer.consume(new GameStarted([{nbFailuresRequiredToFail: 1, nbPeopleOnMission: 3}]));
  expect(receivedReq).to.deep.equal([{nbPeopleOnMission: 3, nbFailuresRequiredToFail: 1}]);
});

test(`LeaderStartedToSelectMembers`, () => {
  let receivedLeader: string
  let teamSelectionStarted = false;
  const gameConsumer = new GameConsumer({assignLeader: leader => {receivedLeader = leader}, startTeamSelection: () => {teamSelectionStarted = true}} as GameStore);
  gameConsumer.consume(new LeaderStartedToSelectMembers("testLeader"));
  expect(teamSelectionStarted).to.be.true;
  expect(receivedLeader).to.equal("testLeader");
});

test(`LeaderSelectedMember`, () => {
  let receivedMember: string
  const gameConsumer = new GameConsumer({selectPlayer: member => {receivedMember = member}} as GameStore);
  gameConsumer.consume(new LeaderSelectedMember("member"));
  expect(receivedMember).to.equal("member");
});

test(`LeaderDeselectedMember`, () => {
  let receivedMember: string
  const gameConsumer = new GameConsumer({deselectPlayer: member => {receivedMember = member}} as GameStore);
  gameConsumer.consume(new LeaderDeselectedMember("member"));
  expect(receivedMember).to.equal("member");
});

test(`LeaderConfirmedTeam`, () => {
  let voteStarted = false;
  const gameConsumer = new GameConsumer({startTeamVote: () => {voteStarted = true}} as GameStore);
  gameConsumer.consume(new LeaderConfirmedTeam());
  expect(voteStarted).to.be.true;
});

test(`PlayerVotedOnTeam`, () => {
  let receivedPlayer: string;
  let receivedApproval: boolean;
  const gameConsumer = new GameConsumer({makePlayerVote: (player, approval) => {
    receivedPlayer = player;
    receivedApproval = approval;
  }} as GameStore);
  gameConsumer.consume(new PlayerVotedOnTeam("testName", true));
  expect(receivedPlayer).to.equal("testName");
  expect(receivedApproval).to.be.true;
});

test(`AllPlayerVotedOnTeam`, () => {
  let receivedApproved: boolean;
  let receivedPlayerVotes: Map<string, boolean>;
  const gameConsumer = new GameConsumer({saveTeamVoteResult: (approved, playerVotes) => {
    receivedApproved = approved;
    receivedPlayerVotes = playerVotes;
  }} as GameStore);
  gameConsumer.consume(new AllPlayerVotedOnTeam(true, new Map<string, boolean>([["p1", true], ["p2", true], ["p3", false]])));
  expect(receivedApproved).to.be.true;
  expect(receivedPlayerVotes).to.deep.equal(new Map<string, boolean>([["p1", true], ["p2", true], ["p3", false]]));
});

test(`MissionStarted`, () => {
  let missionStarted: boolean = false;
  const gameConsumer = new GameConsumer({startMission: () => {missionStarted = true}} as GameStore);
  gameConsumer.consume(new MissionStarted());
  expect(missionStarted).to.be.true;
});

test(`PlayerWorkedOnMission`, () => {
  let receivedPlayer: string;
  let receivedSuccess: boolean;
  const gameConsumer = new GameConsumer({makePlayerWorkOnMission: (player, approval) => {
    receivedPlayer = player;
    receivedSuccess = approval;
  }} as GameStore);
  gameConsumer.consume(new PlayerWorkedOnMission("testName", true));
  expect(receivedPlayer).to.equal("testName");
  expect(receivedSuccess).to.be.true;
  
});

test(`MissionCompleted`, () => {
  let receivedSuccess: boolean;
  let receivedNbFails: number;
  let lastMissionResultShown: boolean;
  const gameConsumer = new GameConsumer({saveMissionResult: (success, nbFails) => {
    receivedSuccess = success;
    receivedNbFails = nbFails;
  },
  showLastMissionResult: ()=>{lastMissionResultShown = true;}
} as GameStore);
  gameConsumer.consume(new MissionCompleted(false, 2));
  expect(receivedSuccess).to.be.false;
  expect(receivedNbFails).to.equal(2);
  expect(lastMissionResultShown).to.be.true;
});

test(`GameEnded`, () => {
  let receivedWinner: Allegiance = null;
  let receivedSpies: Set<string> = null;
  const gameConsumer = new GameConsumer({endGame: (winner, spies) =>{
    receivedWinner = winner;
    receivedSpies = spies;
  }} as GameStore);
  gameConsumer.consume(new GameEnded(Allegiance.Spies, new Set<string>(["a", "b"])));
  expect(receivedWinner).to.equal(Allegiance.Spies);
  expect(receivedSpies).to.deep.equal(new Set<string>(["a", "b"]));
});

