import { expect, test } from "vitest";
import { 
  EventsReplayEnded, 
  EventsReplayStarted, 
  ServerConnectionClosed, 
  ServerConnectionErrorOccured 
} from "../messages/events";
import { ReplayManager, type ReplayStore } from "./replay";

test(`ReplayManager - EventsReplayStarted`, () => {
  let startReplayCalled = false;
  const eventReplayer = new ReplayManager({startReplay: () => {startReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new EventsReplayStarted(null));
  expect(startReplayCalled).to.be.true;
});

test(`ReplayManager - EventsReplayEnded`, () => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new EventsReplayEnded());
  expect(endReplayCalled).to.be.true;
});

test(`ReplayManager - ServerConnectionClosed`, () => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new ServerConnectionClosed());
  expect(endReplayCalled).to.be.true;
});

test(`ReplayManager - ServerConnectionErrorOccured`, () => {
  let endReplayCalled = false;
  const eventReplayer = new ReplayManager({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  eventReplayer.consume(new ServerConnectionErrorOccured());
  expect(endReplayCalled).to.be.true;
});
