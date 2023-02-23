import { expect, test } from "vitest";
import { type Message, MessageBus } from "./message-bus";

test(`Dispatch message to all consumers`, t => {
  const mb = new MessageBus();
  const c1 = []
  const c2 = []
  const c3 = []
  mb.subscribeConsumer({consume: (m:Message) => {
    c1.push(m);
  }})
  mb.subscribeConsumer({consume: (m:Message) => {
    c2.push(m);
  }})
  mb.subscribeConsumer({consume: (m:Message) => {
    c3.push(m);
  }})

  mb.dispatch("m1");
  mb.dispatch("m2");
  mb.dispatch("m3");

  expect(c1).to.deep.equal(["m1", "m2", "m3"]);
  expect(c2).to.deep.equal(["m1", "m2", "m3"]);
  expect(c3).to.deep.equal(["m1", "m2", "m3"]);
});
