import {expect, test} from "vitest";
import { Page } from "../types/types";
import { PageService, type PageValues } from "./Page-service";


test("Is page Lobby", ()=> {
  let service = new PageService({pageToShow: Page.Lobby} as PageValues);
  expect(service.isPageLobby).to.be.true;

  service = new PageService({pageToShow: Page.Game} as PageValues);
  expect(service.isPageLobby).to.be.false;
});

test("Is page Party Room", ()=> {
  let service = new PageService({pageToShow: Page.PartyRoom} as PageValues);
  expect(service.isPagePartyRoom).to.be.true;

  service = new PageService({pageToShow: Page.Game} as PageValues);
  expect(service.isPagePartyRoom).to.be.false;
});

test("Is page Game", ()=> {
  let service = new PageService({pageToShow: Page.Game} as PageValues);
  expect(service.isPageGame).to.be.true;

  service = new PageService({pageToShow: Page.Lobby} as PageValues);
  expect(service.isPageGame).to.be.false;
});