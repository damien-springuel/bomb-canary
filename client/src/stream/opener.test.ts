import test from "ava";
import { AppLoaded, CreatePartySucceeded, JoinPartySucceeded } from "../messages/events";
import { Opener } from "./opener";

test(`Opener - open on AppLoaded`, t => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new AppLoaded())
  t.true(wasCreated);
});

test(`Opener - open on CreatePartySucceeded`, t => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new CreatePartySucceeded())
  t.true(wasCreated);
});

test(`Opener - open on JoinPartySucceded`, t => {
  let wasCreated = false
  const opener = new Opener({create: () => {wasCreated = true;}});
  opener.consume(new JoinPartySucceeded())
  t.true(wasCreated);
});
