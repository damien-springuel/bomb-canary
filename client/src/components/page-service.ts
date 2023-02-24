import { Page } from "../types/types";

export interface PageValues {
  pageToShow: Page
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