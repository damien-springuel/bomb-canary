import type { Dispatcher } from "./messages/dispatcher";
import { AppLoaded } from "./messages/events";
import { Page } from "./store/store";

export interface AppValues {
  pageToShow: Page
}

export class AppService {
  constructor(
    private readonly values: AppValues,
    private readonly dispatcher: Dispatcher){}

  appMounted() {
    this.dispatcher.dispatch(new AppLoaded());
  }

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