import test from "ava";
import { Message, MessageBus } from "./messagebus";

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

  t.deepEqual(c1, ["m1", "m2", "m3"]);
  t.deepEqual(c2, ["m1", "m2", "m3"]);
  t.deepEqual(c3, ["m1", "m2", "m3"]);
});
