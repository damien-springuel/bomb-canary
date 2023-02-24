import { expect, test } from "vitest";
import { IdentityService } from "./Identity-service";


test(`Identity Service - isPlayerIsASpy`, () => {
  let service = new IdentityService({revealedSpies: new Set<string>(["a", "b"]), player: "b"});
  service.isPlayerIsASpy();
  expect(service.isPlayerIsASpy()).to.be.true;
  
  service = new IdentityService({revealedSpies: new Set<string>(), player: "b"});
  expect(service.isPlayerIsASpy()).to.be.false;
});

test(`Identity Service - otherSpies`, () => {
  let service = new IdentityService({revealedSpies: new Set<string>(["a", "b", "c"]), player: "b"});
  service.isPlayerIsASpy();
  expect(service.otherSpies()).to.equal("a, c");
  
  service = new IdentityService({revealedSpies: new Set<string>(), player: "b"});
  expect(service.otherSpies()).to.equal("");
});
