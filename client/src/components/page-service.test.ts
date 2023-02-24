import {expect, test} from "vitest";
import { Page } from "../types/types";
import { PageService } from "./page-service";


test("Is page Lobby", ()=> {
  let service = new PageService({
    pageToShow: Page.Lobby,
  });
  expect(service.isPageLobby).to.be.true;

  service = new PageService({
    pageToShow: Page.Game,
  });
  expect(service.isPageLobby).to.be.false;
});

test("Is page Party Room", ()=> {
  let service = new PageService({
    pageToShow: Page.PartyRoom,
  });
  expect(service.isPagePartyRoom).to.be.true;

  service = new PageService({
    pageToShow: Page.Game,
  });
  expect(service.isPagePartyRoom).to.be.false;
});

test("Is page Party Room", ()=> {
  let service = new PageService({
    pageToShow: Page.Game,
  });
  expect(service.isPageGame).to.be.true;

  service = new PageService({
    pageToShow: Page.Lobby,
  });
  expect(service.isPageGame).to.be.false;
});