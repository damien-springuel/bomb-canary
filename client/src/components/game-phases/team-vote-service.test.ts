import {expect, test} from "vitest";
import { ApproveTeam, RejectTeam } from "../../messages/commands";
import type { Dispatcher } from "../../messages/dispatcher";
import type { Message } from "../../messages/messagebus";
import { TeamVoteService } from "./team-vote-service";

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
  let messageGiven: Message;
  const dispatcher: Dispatcher = {dispatch(message){messageGiven = message}}
  
  const service = new TeamVoteService(
    {
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, dispatcher);

  service.approveTeam();
  expect(messageGiven).to.be.instanceof(ApproveTeam);
  expect(messageGiven).to.deep.equal(new ApproveTeam());
})

test("Reject team", ()=> {
  let messageGiven: Message;
  const dispatcher: Dispatcher = {dispatch(message){messageGiven = message}}
  
  const service = new TeamVoteService(
    {
    player: "a",
    players: ["a", "b", "c"],
    currentTeam: new Set<string>(["a", "b", "c"]),
    peopleThatVotedOnTeam: new Set<string>(["a"]),
    playerVote: true,
  }, dispatcher);

  service.rejectTeam();
  expect(messageGiven).to.be.instanceof(RejectTeam);
  expect(messageGiven).to.deep.equal(new RejectTeam());
})