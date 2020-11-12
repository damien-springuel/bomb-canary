import test from "ava";
import { GameStarted, LeaderConfirmedTeam, LeaderDeselectedMember, LeaderSelectedMember, LeaderStartedToSelectMembers, MissionRequirement, MissionStarted, PlayerVotedOnTeam, PlayerWorkedOnMission } from "../messages/events";
import { GameManager, GameStore } from "./game";

test(`Game Manager - GameStarted`, t => {
  let receivedReq: MissionRequirement[]
  const gameMgr = new GameManager({setMissionRequirements: req => {receivedReq = req}} as GameStore);
  gameMgr.consume(new GameStarted([{nbFailuresRequiredToFail: 1, nbPeopleOnMission: 3}]));
  t.deepEqual(receivedReq, [{nbPeopleOnMission: 3, nbFailuresRequiredToFail: 1}]);
});

test(`Game Manager - LeaderStartedToSelectMembers`, t => {
  let receivedLeader: string
  const gameMgr = new GameManager({assignLeader: leader => {receivedLeader = leader}} as GameStore);
  gameMgr.consume(new LeaderStartedToSelectMembers("testLeader"));
  t.deepEqual(receivedLeader, "testLeader");
});

test(`Game Manager - LeaderSelectedMember`, t => {
  let receivedMember: string
  const gameMgr = new GameManager({selectPlayer: member => {receivedMember = member}} as GameStore);
  gameMgr.consume(new LeaderSelectedMember("member"));
  t.deepEqual(receivedMember, "member");
});

test(`Game Manager - LeaderDeselectedMember`, t => {
  let receivedMember: string
  const gameMgr = new GameManager({deselectPlayer: member => {receivedMember = member}} as GameStore);
  gameMgr.consume(new LeaderDeselectedMember("member"));
  t.deepEqual(receivedMember, "member");
});

test(`Game Manager - LeaderConfirmedTeam`, t => {
  let voteStarted = false;
  const gameMgr = new GameManager({startTeamVote: () => {voteStarted = true}} as GameStore);
  gameMgr.consume(new LeaderConfirmedTeam());
  t.true(voteStarted);
});

test(`Game Manager - PlayerVotedOnTeam`, t => {
  let receivedPlayer: string;
  let receivedApproval: boolean;
  const gameMgr = new GameManager({makePlayerVote: (player, approval) => {
    receivedPlayer = player;
    receivedApproval = approval;
  }} as GameStore);
  gameMgr.consume(new PlayerVotedOnTeam("testName", true));
  t.deepEqual(receivedPlayer, "testName");
  t.deepEqual(receivedApproval, true);
});

test(`Game Manager - MissionStarted`, t => {
  let missionStarted: boolean = false;
  const gameMgr = new GameManager({startMission: () => {missionStarted = true}} as GameStore);
  gameMgr.consume(new MissionStarted());
  t.true(missionStarted);
});


test(`Game Manager - PlayerWorkedOnMission`, t => {
  let receivedPlayer: string;
  let receivedSuccess: boolean;
  const gameMgr = new GameManager({makePlayerWorkOnMission: (player, approval) => {
    receivedPlayer = player;
    receivedSuccess = approval;
  }} as GameStore);
  gameMgr.consume(new PlayerWorkedOnMission("testName", true));
  t.deepEqual(receivedPlayer, "testName");
  t.deepEqual(receivedSuccess, true);
});
