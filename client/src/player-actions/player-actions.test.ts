import test from "ava";
import { HttpPostMock } from "../http/post.test-utils";
import { LeaderConfirmsTeam, LeaderDeselectsMember, LeaderSelectsMember, StartGame } from "../messages/commands";
import { PlayerActions } from "./player-actions";

test(`Start Game`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new StartGame());
  t.deepEqual(httpPost.givenUrl, "/actions/start-game")
});

test(`Leader Selects Member`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderSelectsMember("testName"));
  t.deepEqual(httpPost.givenUrl, "/actions/leader-selects-member");
  t.deepEqual(httpPost.givenData, {member: "testName"});
});

test(`Leader Deselects Member`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderDeselectsMember("testName"));
  t.deepEqual(httpPost.givenUrl, "/actions/leader-deselects-member");
  t.deepEqual(httpPost.givenData, {member: "testName"});
});

test(`Leader Confirms team`, t => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderConfirmsTeam());
  t.deepEqual(httpPost.givenUrl, "/actions/leader-confirms-team");
});
