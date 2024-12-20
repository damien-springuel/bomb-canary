import {expect, test} from "vitest";
import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { TeamSelectionService } from "./TeamSelection-service";

test("Is given player in current team", ()=> {
  const service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
      nbPeopleRequiredOnMission: 3,
    },
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
      nbPeopleRequiredOnMission: 3,
    },
    null);
  expect(service.isPlayerTheLeader).to.be.true;

  service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "b",
      currentTeamVoteNb: 1,
      players: ["a", "b"],
      nbPeopleRequiredOnMission: 3,
    },
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
      nbPeopleRequiredOnMission: 3,
    },
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
      nbPeopleRequiredOnMission: 2,
    },
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
      nbPeopleRequiredOnMission: 3,
    },
    null);
  
  expect(service.canConfirmTeam).to.be.false;

  service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c"],
      nbPeopleRequiredOnMission: 2,
    },
    null);

  expect(service.canConfirmTeam).to.be.true;
});

test("Toggle player selection", ()=> {
  const dispatcher = new DispatcherMock();
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c", "d"],
      nbPeopleRequiredOnMission: 3,
    },
    dispatcher);
  
  service.togglePlayerSelection("d");
  expect(dispatcher.receivedMessage).to.be.an.instanceof(LeaderSelectsMember);
  expect(dispatcher.receivedMessage).to.deep.equal(new LeaderSelectsMember("d"));
  
  dispatcher.receivedMessage = null;
  service.togglePlayerSelection("a");
  expect(dispatcher.receivedMessage).to.be.an.instanceof(LeaderDeselectsMember);
  expect(dispatcher.receivedMessage).to.deep.equal(new LeaderDeselectsMember("a"));
});

test("Confirm Team", ()=> {
  const dispatcher = new DispatcherMock();
  let service = new TeamSelectionService(
    {
      currentTeam: new Set<string>(["a", "b"]),
      player: "a",
      leader: "a",
      currentTeamVoteNb: 1,
      players: ["a", "b", "c", "d"],
      nbPeopleRequiredOnMission: 3,
    },
    dispatcher);
  
  service.confirmTeam();
  expect(dispatcher.receivedMessage).to.be.an.instanceof(LeaderConfirmsTeam);
  expect(dispatcher.receivedMessage).to.deep.equal(new LeaderConfirmsTeam());
});