import test from "ava";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { PartyCreated, PlayerConnected, PlayerDisconnected, PlayerJoined, ServerConnectionClosed, ServerConnectionErrorOccured } from "../messages/events";
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
  handler.onEvent({PlayerJoined: {Name: "testName", Code: "testCode"}});
  t.deepEqual(dispatcher.receivedMessage, new PlayerJoined("testName"));
});
