import { expect, test } from "vitest";
import { 
    AllPlayerVotedOnTeam, 
    GameStarted, 
    LeaderConfirmedTeam, 
    LeaderDeselectedMember, 
    LeaderSelectedMember, 
    LeaderStartedToSelectMembers, 
    MissionCompleted, 
    type MissionRequirement, 
    MissionStarted, 
    PlayerVotedOnTeam, 
    PlayerWorkedOnMission 
} from "../messages/events";
import { GameManager, type GameStore } from "./game";

test(`Game Manager - GameStarted`, t => {
  let receivedReq: MissionRequirement[]
  const gameMgr = new GameManager({setMissionRequirements: req => {receivedReq = req}} as GameStore);
  gameMgr.consume(new GameStarted([{nbFailuresRequiredToFail: 1, nbPeopleOnMission: 3}]));
  expect(receivedReq).to.deep.equal([{nbPeopleOnMission: 3, nbFailuresRequiredToFail: 1}]);
});

test(`Game Manager - LeaderStartedToSelectMembers`, t => {
  let receivedLeader: string
  let teamSelectionStarted = false;
  const gameMgr = new GameManager({assignLeader: leader => {receivedLeader = leader}, startTeamSelection: () => {teamSelectionStarted = true}} as GameStore);
  gameMgr.consume(new LeaderStartedToSelectMembers("testLeader"));
  expect(teamSelectionStarted).to.be.true;
  expect(receivedLeader).to.equal("testLeader");
});

test(`Game Manager - LeaderSelectedMember`, t => {
  let receivedMember: string
  const gameMgr = new GameManager({selectPlayer: member => {receivedMember = member}} as GameStore);
  gameMgr.consume(new LeaderSelectedMember("member"));
  expect(receivedMember).to.equal("member");
});

test(`Game Manager - LeaderDeselectedMember`, t => {
  let receivedMember: string
  const gameMgr = new GameManager({deselectPlayer: member => {receivedMember = member}} as GameStore);
  gameMgr.consume(new LeaderDeselectedMember("member"));
  expect(receivedMember).to.equal("member");
});

test(`Game Manager - LeaderConfirmedTeam`, t => {
  let voteStarted = false;
  const gameMgr = new GameManager({startTeamVote: () => {voteStarted = true}} as GameStore);
  gameMgr.consume(new LeaderConfirmedTeam());
  expect(voteStarted).to.be.true;
});

test(`Game Manager - PlayerVotedOnTeam`, t => {
  let receivedPlayer: string;
  let receivedApproval: boolean;
  const gameMgr = new GameManager({makePlayerVote: (player, approval) => {
    receivedPlayer = player;
    receivedApproval = approval;
  }} as GameStore);
  gameMgr.consume(new PlayerVotedOnTeam("testName", true));
  expect(receivedPlayer).to.equal("testName");
  expect(receivedApproval).to.be.true;
});

test(`Game Manager - AllPlayerVotedOnTeam`, t => {
  let receivedApproved: boolean;
  let receivedPlayerVotes: Map<string, boolean>;
  const gameMgr = new GameManager({saveTeamVoteResult: (approved, playerVotes) => {
    receivedApproved = approved;
    receivedPlayerVotes = playerVotes;
  }} as GameStore);
  gameMgr.consume(new AllPlayerVotedOnTeam(true, new Map<string, boolean>([["p1", true], ["p2", true], ["p3", false]])));
  expect(receivedApproved).to.be.true;
  expect(receivedPlayerVotes).to.deep.equal(new Map<string, boolean>([["p1", true], ["p2", true], ["p3", false]]));
});

test(`Game Manager - MissionStarted`, t => {
  let missionStarted: boolean = false;
  const gameMgr = new GameManager({startMission: () => {missionStarted = true}} as GameStore);
  gameMgr.consume(new MissionStarted());
  expect(missionStarted).to.be.true;
});

test(`Game Manager - PlayerWorkedOnMission`, t => {
  let receivedPlayer: string;
  let receivedSuccess: boolean;
  const gameMgr = new GameManager({makePlayerWorkOnMission: (player, approval) => {
    receivedPlayer = player;
    receivedSuccess = approval;
  }} as GameStore);
  gameMgr.consume(new PlayerWorkedOnMission("testName", true));
  expect(receivedPlayer).to.equal("testName");
  expect(receivedSuccess).to.be.true;
  
});

test(`Game Manager - SaveMissionResult`, t => {
  let receivedSuccess: boolean;
  let receivedNbFails: number;
  const gameMgr = new GameManager({saveMissionResult: (success, nbFails) => {
    receivedSuccess = success;
    receivedNbFails = nbFails;
  }} as GameStore);
  gameMgr.consume(new MissionCompleted(false, 2));
  expect(receivedSuccess).to.be.false;
  expect(receivedNbFails).to.equal(2);
});

