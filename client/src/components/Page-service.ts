import { Page } from "../types/types";
import type { GameValues } from "./Game-service";
import type { PartyRoomValues } from "./PartyRoom-service";

export interface PageValues {
  readonly pageToShow: Page,
  readonly gameValues: GameValues,
  readonly partyRoomValues: PartyRoomValues,
}

export class PageService {
  constructor(private readonly values: PageValues){}

  private isPage(page: Page) {
    return this.values.pageToShow == page;
  }

  get isPageLobby(): boolean {
    return this.isPage(Page.Lobby);
  }

  get isPagePartyRoom(): boolean {
    return this.isPage(Page.PartyRoom);
  }

  get isPageGame(): boolean {
    return this.isPage(Page.Game);
  }
}