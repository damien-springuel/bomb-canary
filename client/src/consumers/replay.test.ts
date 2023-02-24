import { expect, test } from "vitest";
import { 
  EventsReplayEnded, 
  EventsReplayStarted, 
  ServerConnectionClosed, 
  ServerConnectionErrorOccured 
} from "../messages/events";
import { ReplayConsumer, type ReplayStore } from "./replay";

test(`ReplayManager - EventsReplayStarted`, () => {
  let startReplayCalled = false;
  const replayConsumer = new ReplayConsumer({startReplay: () => {startReplayCalled = true;}} as ReplayStore);
  replayConsumer.consume(new EventsReplayStarted(null));
  expect(startReplayCalled).to.be.true;
});

test(`ReplayManager - EventsReplayEnded`, () => {
  let endReplayCalled = false;
  const replayConsumer = new ReplayConsumer({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  replayConsumer.consume(new EventsReplayEnded());
  expect(endReplayCalled).to.be.true;
});

test(`ReplayManager - ServerConnectionClosed`, () => {
  let endReplayCalled = false;
  const replayConsumer = new ReplayConsumer({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  replayConsumer.consume(new ServerConnectionClosed());
  expect(endReplayCalled).to.be.true;
});

test(`ReplayManager - ServerConnectionErrorOccured`, () => {
  let endReplayCalled = false;
  const replayConsumer = new ReplayConsumer({endReplay: () => {endReplayCalled = true;}} as ReplayStore);
  replayConsumer.consume(new ServerConnectionErrorOccured());
  expect(endReplayCalled).to.be.true;
});
