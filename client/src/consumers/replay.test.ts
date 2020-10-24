import test from "ava";
import { EventsReplayEnded, ServerConnectionClosed, ServerConnectionErrorOccured } from "../messages/events";
import { ReplayManager } from "./replay";

test(`ReplayManager - EventsReplayEnded`, t => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}});
  eventReplayer.consume(new EventsReplayEnded());
  t.true(endReplayCalled);
});

test(`ReplayManager - ServerConnectionClosed`, t => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}});
  eventReplayer.consume(new ServerConnectionClosed());
  t.true(endReplayCalled);
});

test(`ReplayManager - ServerConnectionErrorOccured`, t => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}});
  eventReplayer.consume(new ServerConnectionErrorOccured());
  t.true(endReplayCalled);
});
