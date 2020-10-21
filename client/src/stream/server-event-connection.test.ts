import test from "ava";
import { AppLoaded, PartyCreated } from "../messages/events";
import { ServerEventConnectionOpener } from "./server-event-connection";

test(`Server Event Connection - open on AppLoaded`, t => {
  let wasOpen = false
  const opener = {open: () => {wasOpen = true;}};
  const serverEventConnectionOpener = new ServerEventConnectionOpener(opener);
  serverEventConnectionOpener.consume(new AppLoaded())
  t.true(wasOpen);
});

test(`Server Event Connection - open on PartyCreated`, t => {
  let wasOpen = false
  const opener = {open: () => {wasOpen = true;}};
  const serverEventConnectionOpener = new ServerEventConnectionOpener(opener);
  serverEventConnectionOpener.consume(new PartyCreated(null))
  t.true(wasOpen);
});
