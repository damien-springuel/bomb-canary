import {expect, test} from "vitest";
import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember } from "../../messages/commands";
import type { Dispatcher } from "../../messages/dispatcher";
import type { Message } from "../../messages/messagebus";
import { MissionTrackerService } from "../mission-tracker-service";
import { TeamSelectionService } from "./team-selection-service";

test("Is given player in current team", ()=> {
  const service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
      
    },
    new MissionTrackerService({missionResults: [], missionRequirements:[]}),
    null);

  expect(service.isGivenPlayerInTeam("a")).to.be.true;
  expect(service.isGivenPlayerInTeam("c")).to.be.false;
});

test("Is the player the leader", ()=> {
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
    },
    new MissionTrackerService({missionResults: [], missionRequirements:[]}),
    null);
  expect(service.isPlayerTheLeader).to.be.true;

  service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "b",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
    },
    new MissionTrackerService({missionResults: [], missionRequirements:[]}),
    null);

  expect(service.isPlayerTheLeader).to.be.false;
});

test("Is the player selectable to the team - team not full", ()=> {
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c"],
    },
    new MissionTrackerService({
      missionResults: [], 
      missionRequirements:[
        {nbPeopleOnMission: 3, nbFailuresRequiredToFail: 1}
      ]}),
    null);
  expect(service.isGivenPlayerSelectableForTeam("a")).to.be.true;
  expect(service.isGivenPlayerSelectableForTeam("b")).to.be.true;
  expect(service.isGivenPlayerSelectableForTeam("c")).to.be.true;
});

test("Is the player selectable to the team - team full", ()=> {
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c"],
    },
    new MissionTrackerService({
      missionResults: [], 
      missionRequirements:[
        {nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1}
      ]}),
    null);
  expect(service.isGivenPlayerSelectableForTeam("a")).to.be.true;
  expect(service.isGivenPlayerSelectableForTeam("b")).to.be.true;
  expect(service.isGivenPlayerSelectableForTeam("c")).to.be.false;
});

test("Can confirm team", ()=> {
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c"],
    },
    new MissionTrackerService({
      missionResults: [], 
      missionRequirements:[
        {nbPeopleOnMission: 3, nbFailuresRequiredToFail: 1}
      ]}),
    null);
  
  expect(service.canConfirmTeam).to.be.false;

  service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c"],
    },
    new MissionTrackerService({
      missionResults: [], 
      missionRequirements:[
        {nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1}
      ]}),
    null);

  expect(service.canConfirmTeam).to.be.true;
});

test("Toggle player selection", ()=> {
  let messageGiven: Message;
  const dispatcher: Dispatcher = {dispatch(message){messageGiven = message}};
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c", "d"],
    },
    new MissionTrackerService({
      missionResults: [], 
      missionRequirements:[
        {nbPeopleOnMission: 3, nbFailuresRequiredToFail: 1}
      ]}),
    dispatcher);
  
  service.togglePlayerSelection("d");
  console.log(messageGiven);
  expect(messageGiven).to.be.an.instanceof(LeaderSelectsMember);
  expect(messageGiven).to.deep.equal(new LeaderSelectsMember("d"));
  
  messageGiven = null;
  service.togglePlayerSelection("a");
  expect(messageGiven).to.be.an.instanceof(LeaderDeselectsMember);
  expect(messageGiven).to.deep.equal(new LeaderDeselectsMember("a"));
});

test("Confirm Team", ()=> {
  let messageGiven: Message;
  const dispatcher: Dispatcher = {dispatch(message){messageGiven = message}};
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c", "d"],
    },
    new MissionTrackerService({
      missionResults: [], 
      missionRequirements:[
        {nbPeopleOnMission: 2, nbFailuresRequiredToFail: 1}
      ]}),
    dispatcher);
  
  service.confirmTeam();
  expect(messageGiven).to.be.an.instanceof(LeaderConfirmsTeam);
  expect(messageGiven).to.deep.equal(new LeaderConfirmsTeam());
});