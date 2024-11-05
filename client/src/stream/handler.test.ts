import { expect, test } from "vitest";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { 
  AllPlayerVotedOnTeam, 
  EventsReplayEnded, 
  EventsReplayStarted, 
  GameStarted, 
  LeaderConfirmedTeam, 
  LeaderDeselectedMember, 
  LeaderSelectedMember, 
  LeaderStartedToSelectMembers, 
  MissionCompleted, 
  MissionStarted, 
  PlayerConnected, 
  PlayerDisconnected, 
  PlayerJoined, 
  PlayerVotedOnTeam, 
  PlayerWorkedOnMission, 
  ServerConnectionClosed, 
  ServerConnectionErrorOccured, 
  SpiesRevealed 
} from "../messages/events";
import { Handler } from "./handler";

test(`Handler - onClose`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onClose();
  expect(dispatcher.receivedMessage).to.deep.equal(new ServerConnectionClosed());
});

test(`Handler - onError`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onError();
  expect(dispatcher.receivedMessage).to.deep.equal(new ServerConnectionErrorOccured());
});

test(`Handler - onEvent - EventsReplayStarted`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({EventsReplayStarted: {Player: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new EventsReplayStarted("testName"));
});

test(`Handler - onEvent - EventsReplayEnded`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({EventsReplayEnded: {}});
  expect(dispatcher.receivedMessage).to.deep.equal(new EventsReplayEnded());
});

test(`Handler - onEvent - PlayerConnected`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerConnected: {Name: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new PlayerConnected("testName"));
});

test(`Handler - onEvent - PlayerDisconnected`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerDisconnected: {Name: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new PlayerDisconnected("testName"));
});

test(`Handler - onEvent - PlayerJoined`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerJoined: {Name: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new PlayerJoined("testName"));
});

test(`Handler - onEvent - GameStarted`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({GameStarted: {MissionRequirements: [{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 2}]}});
  expect(dispatcher.receivedMessage).to.deep.equal(new GameStarted([{nbPeopleOnMission: 3, nbFailuresRequiredToFail: 2}]));
});

test(`Handler - onEvent - SpiesRevealed`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({SpiesRevealed: {Spies: {"name1": {}, "name2": {}}}});
  expect(dispatcher.receivedMessage).to.deep.equal(new SpiesRevealed(new Set<string>(["name1", "name2"])));
});

test(`Handler - onEvent - SpiesRevealed - no spies`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({SpiesRevealed: {Spies: {}}});
  expect(dispatcher.receivedMessage).to.deep.equal(new SpiesRevealed(new Set<string>([])));
});

test(`Handler - onEvent - LeaderStartedToSelectMembers`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderStartedToSelectMembers: {Leader: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new LeaderStartedToSelectMembers("testName"));
});

test(`Handler - onEvent - LeaderSelectedMember`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderSelectedMember: {SelectedMember: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new LeaderSelectedMember("testName"));
});

test(`Handler - onEvent - LeaderDeselectedMember`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderDeselectedMember: {DeselectedMember: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new LeaderDeselectedMember("testName"));
});

test(`Handler - onEvent - LeaderConfirmedTeam`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderConfirmedSelection: {}});
  expect(dispatcher.receivedMessage).to.deep.equal(new LeaderConfirmedTeam());
});

test(`Handler - onEvent - PlayerVotedOnTeam - with 'Approved' field`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerVotedOnTeam: {Player: "testName", Approved: false}});
  expect(dispatcher.receivedMessage).to.deep.equal(new PlayerVotedOnTeam("testName", false));
});

test(`Handler - onEvent - PlayerVotedOnTeam - without 'Approved' field`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerVotedOnTeam: {Player: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new PlayerVotedOnTeam("testName", null));
});

test(`Handler - onEvent - AllPlayerVoted`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({AllPlayerVotedOnTeam: {Approved: true, VoteFailures: 3, PlayerVotes: {"Alice": true, "Bob": false}}});
  expect(dispatcher.receivedMessage).to.deep.equal(new AllPlayerVotedOnTeam(true, new Map<string, boolean>([["Alice", true], ["Bob", false]])));
});

test(`Handler - onEvent - MissionStarted`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({MissionStarted: {}});
  expect(dispatcher.receivedMessage).to.deep.equal(new MissionStarted());
});

test(`Handler - onEvent - PlayerWorkedOnMission - with success field`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerWorkedOnMission: {Player: "testName", Success: false}});
  expect(dispatcher.receivedMessage).to.deep.equal(new PlayerWorkedOnMission("testName", false));
});

test(`Handler - onEvent - PlayerWorkedOnMission - without success field`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerWorkedOnMission: {Player: "testName"}});
  expect(dispatcher.receivedMessage).to.deep.equal(new PlayerWorkedOnMission("testName", null));
});

test(`Handler - onEvent - MissionCompleted`, () => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({MissionCompleted: {Success: false, NbFails: 1}});
  expect(dispatcher.receivedMessage).to.deep.equal(new MissionCompleted(false, 1));
});
