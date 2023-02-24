import { expect, test } from "vitest";
import { HttpPostMock } from "../http/post.test-utils";
import { 
  ApproveTeam, 
  FailMission, 
  LeaderConfirmsTeam, 
  LeaderDeselectsMember, 
  LeaderSelectsMember, 
  RejectTeam, 
  StartGame, 
  SucceedMission 
} from "../messages/commands";
import { PlayerActions } from "./player-actions";

test(`Player Actions - Start Game`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new StartGame());
  expect(httpPost.givenUrl).to.equal("/actions/start-game");
});

test(`Player Actions - Leader Selects Member`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderSelectsMember("testName"));
  expect(httpPost.givenUrl).to.equal("/actions/leader-selects-member");
  expect(httpPost.givenData).to.deep.equal({member: "testName"});
});

test(`Player Actions - Leader Deselects Member`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderDeselectsMember("testName"));
  expect(httpPost.givenUrl).to.equal("/actions/leader-deselects-member");
  expect(httpPost.givenData).to.deep.equal({member: "testName"});
});

test(`Player Actions - Leader Confirms team`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new LeaderConfirmsTeam());
  expect(httpPost.givenUrl).to.equal("/actions/leader-confirms-team");
});

test(`Player Actions - Approve Team`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new ApproveTeam());
  expect(httpPost.givenUrl).to.equal("/actions/approve-team");
});

test(`Player Actions - Reject Team`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new RejectTeam());
  expect(httpPost.givenUrl).to.equal("/actions/reject-team");
});

test(`Player Actions - Succeed Mission`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new SucceedMission());
  expect(httpPost.givenUrl).to.equal("/actions/succeed-mission");
});

test(`Player Actions - Fail Mission`, () => {
  const httpPost = new HttpPostMock();
  const playerActions = new PlayerActions(httpPost);
  playerActions.consume(new FailMission());
  expect(httpPost.givenUrl).to.equal("/actions/fail-mission");
});
