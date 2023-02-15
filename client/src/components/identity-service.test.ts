import { expect, test } from "vitest";
import type { StoreValues } from "../store/store";
import { IdentityService } from "./identity-service";


test(`Identity Service - isPlayerIsASpy`, t => {
  const storeValues: StoreValues = {revealedSpies: new Set<string>(["a", "b"]), player: "b"} as StoreValues
  const service = new IdentityService(storeValues);
  service.isPlayerIsASpy();
  expect(service.isPlayerIsASpy()).to.be.true;

  storeValues.revealedSpies = new Set<string>();
  expect(service.isPlayerIsASpy()).to.be.false;
});

test(`Identity Service - otherSpies`, t => {
  const storeValues: StoreValues = {revealedSpies: new Set<string>(["a", "b", "c"]), player: "b"} as StoreValues
  const service = new IdentityService(storeValues);
  expect(service.otherSpies()).to.deep.equal(["a", "c"]);
  
  storeValues.revealedSpies = new Set<string>();
  expect(service.otherSpies()).to.deep.equal([]);
});
