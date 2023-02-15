import { expect, test } from "vitest";
import { AppLoaded, CreatePartySucceeded, JoinPartySucceeded } from "../messages/events";
import { Opener } from "./opener";

test(`Opener - open on AppLoaded`, t => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new AppLoaded())
  expect(wasCreated).to.be.true;
});

test(`Opener - open on CreatePartySucceeded`, t => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new CreatePartySucceeded())
  expect(wasCreated).to.be.true;
});

test(`Opener - open on JoinPartySucceded`, t => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new JoinPartySucceeded())
  expect(wasCreated).to.be.true;
});
