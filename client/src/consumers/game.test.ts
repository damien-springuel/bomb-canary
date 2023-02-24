import { expect, test } from "vitest";
import { 
    AllPlayerVotedOnTeam, 
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
import type { MissionRequirement } from "../types/types";
import { GameConsumer, type GameStore } from "./game";

test(`Game Manager - GameStarted`, t => {
  let receivedReq: MissionRequirement[]
  const gameConsumer = new GameConsumer({setMissionRequirements: req => {receivedReq = req}} as GameStore);
  gameConsumer.consume(new GameStarted([{nbFailuresRequiredToFail: 1, nbPeopleOnMission: 3}]));
  expect(receivedReq).to.deep.equal([{nbPeopleOnMission: 3, nbFailuresRequiredToFail: 1}]);
});

test(`Game Manager - LeaderStartedToSelectMembers`, t => {
  let receivedLeader: string
  let teamSelectionStarted = false;
  const gameConsumer = new GameConsumer({assignLeader: leader => {receivedLeader = leader}, startTeamSelection: () => {teamSelectionStarted = true}} as GameStore);
  gameConsumer.consume(new LeaderStartedToSelectMembers("testLeader"));
  expect(teamSelectionStarted).to.be.true;
  expect(receivedLeader).to.equal("testLeader");
});

test(`Game Manager - LeaderSelectedMember`, t => {
  let receivedMember: string
  const gameConsumer = new GameConsumer({selectPlayer: member => {receivedMember = member}} as GameStore);
  gameConsumer.consume(new LeaderSelectedMember("member"));
  expect(receivedMember).to.equal("member");
});

test(`Game Manager - LeaderDeselectedMember`, t => {
  let receivedMember: string
  const gameConsumer = new GameConsumer({deselectPlayer: member => {receivedMember = member}} as GameStore);
  gameConsumer.consume(new LeaderDeselectedMember("member"));
  expect(receivedMember).to.equal("member");
});

test(`Game Manager - LeaderConfirmedTeam`, t => {
  let voteStarted = false;
  const gameConsumer = new GameConsumer({startTeamVote: () => {voteStarted = true}} as GameStore);
  gameConsumer.consume(new LeaderConfirmedTeam());
  expect(voteStarted).to.be.true;
});

test(`Game Manager - PlayerVotedOnTeam`, t => {
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

test(`Game Manager - AllPlayerVotedOnTeam`, t => {
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

test(`Game Manager - MissionStarted`, t => {
  let missionStarted: boolean = false;
  const gameConsumer = new GameConsumer({startMission: () => {missionStarted = true}} as GameStore);
  gameConsumer.consume(new MissionStarted());
  expect(missionStarted).to.be.true;
});

test(`Game Manager - PlayerWorkedOnMission`, t => {
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

test(`Game Manager - SaveMissionResult`, t => {
  let receivedSuccess: boolean;
  let receivedNbFails: number;
  const gameConsumer = new GameConsumer({saveMissionResult: (success, nbFails) => {
    receivedSuccess = success;
    receivedNbFails = nbFails;
  }} as GameStore);
  gameConsumer.consume(new MissionCompleted(false, 2));
  expect(receivedSuccess).to.be.false;
  expect(receivedNbFails).to.equal(2);
});

