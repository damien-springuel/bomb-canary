import {expect, test} from "vitest";
import { ApproveTeam, RejectTeam } from "../messages/commands";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { TeamVoteService } from "./TeamVote-service";

test("Current Team", ()=> {
  const service = new TeamVoteService({
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, null);

  expect(service.currentTeamAsString).to.equal("a, b and c");
});

test("Has given player voted", ()=> {
  const service = new TeamVoteService({
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, null);

  expect(service.hasGivenPlayerVoted("a")).to.true;
  expect(service.hasGivenPlayerVoted("b")).to.false;
})

test("Has current player voted", ()=> {
  const service = new TeamVoteService({
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, null);

  expect(service.hasCurrentPlayerVoted).to.true;
})

test("Get player vote", ()=> {
  const service = new TeamVoteService({
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, null);

  expect(service.playerVote).to.true;
})

test("Get players", ()=> {
  const service = new TeamVoteService({
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, null);

  expect(service.players).to.deep.equal(["a", "b", "c"]);
})

test("Approve team", ()=> {
  const dispatcher = new DispatcherMock();
  
  const service = new TeamVoteService(
    {
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, dispatcher);

  service.approveTeam();
  expect(dispatcher.receivedMessage).to.be.instanceof(ApproveTeam);
  expect(dispatcher.receivedMessage).to.deep.equal(new ApproveTeam());
})

test("Reject team", ()=> {
  const dispatcher = new DispatcherMock();
  
  const service = new TeamVoteService(
    {
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, dispatcher);

  service.rejectTeam();
  expect(dispatcher.receivedMessage).to.be.instanceof(RejectTeam);
  expect(dispatcher.receivedMessage).to.deep.equal(new RejectTeam());
})