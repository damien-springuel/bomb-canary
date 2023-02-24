import { expect, test } from "vitest";
import { ServerConnectionClosed } from "../messages/events";
import { ResetConsumer } from "./reset";

test(`Reset manager - ServerConnectionClosed `, () => {
  let reset = false;
  const resetConsumer = new ResetConsumer({reset: () => {reset = true;}});
  resetConsumer.consume(new ServerConnectionClosed());
  expect(reset).to.be.true;
});
