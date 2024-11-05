import { expect, test } from "vitest";
import { AppLoaded, JoinPartySucceeded } from "../messages/events";
import { Opener } from "./opener";

test(`Opener - open on AppLoaded`, () => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new AppLoaded())
  expect(wasCreated).to.be.true;
});

test(`Opener - open on JoinPartySucceded`, () => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new JoinPartySucceeded())
  expect(wasCreated).to.be.true;
});
