import { expect, test } from "vitest";
import { ServerConnectionClosed } from "../messages/events";
import { ResetManager } from "./reset";

test(`Reset manager - ServerConnectionClosed `, t => {
  let reset = false;
  const eventReplayer = new ResetManager({reset: () => {reset = true;}});
  eventReplayer.consume(new ServerConnectionClosed());
  expect(reset).to.be.true;
});
