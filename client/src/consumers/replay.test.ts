import test from "ava";
import { EventsReplayEnded, EventsReplayStarted, ServerConnectionClosed, ServerConnectionErrorOccured } from "../messages/events";
import { ReplayManager, ReplayStore } from "./replay";

test(`ReplayManager - EventsReplayStarted`, t => {
  let startReplayCalled = false;
  const eventReplayer = new ReplayManager({startReplay: () => {startReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new EventsReplayStarted(null));
  t.true(startReplayCalled);
});

test(`ReplayManager - EventsReplayEnded`, t => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new EventsReplayEnded());
  t.true(endReplayCalled);
});

test(`ReplayManager - ServerConnectionClosed`, t => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new ServerConnectionClosed());
  t.true(endReplayCalled);
});

test(`ReplayManager - ServerConnectionErrorOccured`, t => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new ServerConnectionErrorOccured());
  t.true(endReplayCalled);
});
