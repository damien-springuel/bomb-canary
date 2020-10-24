import test from "ava";
import { DispatcherMock } from "../messages/dispatcher.test-utils";
import { PlayerConnected, PlayerDisconnected, PlayerJoined, ServerConnectionClosed, ServerConnectionErrorOccured } from "../messages/events";
import { Handler } from "./handler";
import type { ServerEvent } from "./server-event";

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

test(`Handler - onEvent - PlayerConnected`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerConnected: {Name: "testName"}} as ServerEvent);
  t.deepEqual(dispatcher.receivedMessage, new PlayerConnected("testName"));
});

test(`Handler - onEvent - PlayerDisconnected`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerDisconnected: {Name: "testName"}} as ServerEvent);
  t.deepEqual(dispatcher.receivedMessage, new PlayerDisconnected("testName"));
});

test(`Handler - onEvent - PlayerJoined`, t => {
  const dispatcher: DispatcherMock = new DispatcherMock();
  const handler = new Handler(dispatcher);
  handler.onEvent({PlayerJoined: {Name: "testName", Code: "testCode"}} as ServerEvent);
  t.deepEqual(dispatcher.receivedMessage, new PlayerJoined("testName", "testCode"));
});
