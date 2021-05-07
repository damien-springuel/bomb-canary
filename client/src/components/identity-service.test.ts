import test from "ava";
import { CloseIdentity } from "../messages/commands";
import type { Message } from "../messages/messagebus";
import type { StoreValues } from "../store/store";
import { IdentityService } from "./identity-service";

test(`Identity Service - close identity`, t => {
  let lastMessage: Message;
  const dispatcher = {dispatch: (m: Message) => {lastMessage = m}};
  const service = new IdentityService(dispatcher, null);
  service.closeIdentity();
  t.deepEqual(lastMessage, new CloseIdentity());
});

test(`Identity Service - isPlayerIsASpy`, t => {
  const storeValues: StoreValues = {revealedSpies: new Set<string>(["a", "b"]), player: "b"} as StoreValues
  const service = new IdentityService(null, storeValues);
  service.isPlayerIsASpy();
  t.true(service.isPlayerIsASpy());

  storeValues.revealedSpies = new Set<string>();
  t.false(service.isPlayerIsASpy());
});

test(`Identity Service - otherSpies`, t => {
  const storeValues: StoreValues = {revealedSpies: new Set<string>(["a", "b", "c"]), player: "b"} as StoreValues
  const service = new IdentityService(null, storeValues);
  t.deepEqual(service.otherSpies(), ["a", "c"]);
  
  storeValues.revealedSpies = new Set<string>();
  t.deepEqual(service.otherSpies(), []);
});
