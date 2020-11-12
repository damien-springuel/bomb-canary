import test from "ava";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { AllPlayerVotedOnTeam, EventsReplayEnded, EventsReplayStarted, GameStarted, LeaderConfirmedTeam, LeaderDeselectedMember, LeaderSelectedMember, LeaderStartedToSelectMembers, MissionCompleted, MissionStarted, PartyCreated, PlayerConnected, PlayerDisconnected, PlayerJoined, PlayerVotedOnTeam, PlayerWorkedOnMission, ServerConnectionClosed, ServerConnectionErrorOccured, SpiesRevealed } from "../messages/events";
import { Handler } from "./handler";

test(`Handler - onClose`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onClose();
  t.deepEqual(dispatcher.receivedMessage, new ServerConnectionClosed());
});

test(`Handler - onError`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onError();
  t.deepEqual(dispatcher.receivedMessage, new ServerConnectionErrorOccured());
});

test(`Handler - onEvent - EventsReplayStarted`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({EventsReplayStarted: {Player: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new EventsReplayStarted("testName"));
});

test(`Handler - onEvent - EventsReplayEnded`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({EventsReplayEnded: {}});
  t.deepEqual(dispatcher.receivedMessage, new EventsReplayEnded());
});

test(`Handler - onEvent - PartyCreated`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PartyCreated: {Code: "testCode"}});
  t.deepEqual(dispatcher.receivedMessage, new PartyCreated("testCode"));
});

test(`Handler - onEvent - PlayerConnected`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerConnected: {Name: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerConnected("testName"));
});

test(`Handler - onEvent - PlayerDisconnected`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerDisconnected: {Name: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerDisconnected("testName"));
});

test(`Handler - onEvent - PlayerJoined`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerJoined: {Name: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerJoined("testName"));
});

test(`Handler - onEvent - GameStarted`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({GameStarted: {MissionRequirements: [{NbPeopleOnMission: 3, NbFailuresRequiredToFail: 2}]}});
  t.deepEqual(dispatcher.receivedMessage, new GameStarted([{nbPeopleOnMission: 3, nbFailuresRequiredToFail: 2}]));
});

test(`Handler - onEvent - SpiesRevealed`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({SpiesRevealed: {Spies: {"name1": {}, "name2": {}}}});
  t.deepEqual(dispatcher.receivedMessage, new SpiesRevealed(new Set<string>(["name1", "name2"])));
});

test(`Handler - onEvent - SpiesRevealed - no spies`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({SpiesRevealed: {Spies: {}}});
  t.deepEqual(dispatcher.receivedMessage, new SpiesRevealed(new Set<string>([])));
});

test(`Handler - onEvent - LeaderStartedToSelectMembers`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderStartedToSelectMembers: {Leader: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new LeaderStartedToSelectMembers("testName"));
});

test(`Handler - onEvent - LeaderSelectedMember`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderSelectedMember: {SelectedMember: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new LeaderSelectedMember("testName"));
});

test(`Handler - onEvent - LeaderDeselectedMember`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderDeselectedMember: {DeselectedMember: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new LeaderDeselectedMember("testName"));
});

test(`Handler - onEvent - LeaderConfirmedTeam`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({LeaderConfirmedSelection: {}});
  t.deepEqual(dispatcher.receivedMessage, new LeaderConfirmedTeam());
});

test(`Handler - onEvent - PlayerVotedOnTeam - with 'Approved' field`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerVotedOnTeam: {Player: "testName", Approved: false}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerVotedOnTeam("testName", false));
});

test(`Handler - onEvent - PlayerVotedOnTeam - without 'Approved' field`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerVotedOnTeam: {Player: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerVotedOnTeam("testName", null));
});

test(`Handler - onEvent - AllPlayerVoted`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({AllPlayerVotedOnTeam: {Approved: true, VoteFailures: 3, PlayerVotes: {"Alice": true, "Bob": false}}});
  t.deepEqual(dispatcher.receivedMessage, new AllPlayerVotedOnTeam(true, 3, new Map<string, boolean>([["Alice", true], ["Bob", false]])));
});

test(`Handler - onEvent - MissionStarted`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({MissionStarted: {}});
  t.deepEqual(dispatcher.receivedMessage, new MissionStarted());
});

test(`Handler - onEvent - PlayerWorkedOnMission - with success field`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerWorkedOnMission: {Player: "testName", Success: false}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerWorkedOnMission("testName", false));
});

test(`Handler - onEvent - PlayerWorkedOnMission - without success field`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerWorkedOnMission: {Player: "testName"}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerWorkedOnMission("testName", null));
});

test(`Handler - onEvent - MissionCompleted`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({MissionCompleted: {Success: false, NbFails: 1}});
  t.deepEqual(dispatcher.receivedMessage, new MissionCompleted(false, 1));
});
