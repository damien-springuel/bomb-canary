import test from "ava";
import { HttpPostMock } from "../http/post.test-utils";
import { ApproveTeam, LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember, RejectTeam, StartGame } from "../messages/commands";
import { PlayerActions } from "./player-actions";

test(`Player Actions - Start Game`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new StartGame());
  t.deepEqual(httpPost.givenUrl, "/actions/start-game")
});

test(`Player Actions - Leader Selects Member`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderSelectsMember("testName"));
  t.deepEqual(httpPost.givenUrl, "/actions/leader-selects-member");
  t.deepEqual(httpPost.givenData, {member: "testName"});
});

test(`Player Actions - Leader Deselects Member`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderDeselectsMember("testName"));
  t.deepEqual(httpPost.givenUrl, "/actions/leader-deselects-member");
  t.deepEqual(httpPost.givenData, {member: "testName"});
});

test(`Player Actions - Leader Confirms team`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderConfirmsTeam());
  t.deepEqual(httpPost.givenUrl, "/actions/leader-confirms-team");
});

test(`Player Actions - Approve Team`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new ApproveTeam());
  t.deepEqual(httpPost.givenUrl, "/actions/approve-team");
});

test(`Player Actions - Reject Team`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new RejectTeam());
  t.deepEqual(httpPost.givenUrl, "/actions/reject-team");
});
