import type { Dispatcher } from "./messages/dispatcher";
import { AppLoaded } from "./messages/events";
import { Page } from "./types/types";

export interface AppValues {
  readonly pageToShow: Page,
}

export class AppService {
  constructor(
    private readonly dispatcher: Dispatcher,
    private readonly values: AppValues,
  ){}

  appMounted() {
    this.dispatcher.dispatch(new AppLoaded());
  }

  private isPage(page: Page) {
    return this.values.pageToShow == page;
  }

  get isPagePartyRoom(): boolean {
    return this.isPage(Page.PartyRoom);
  }

  get isPageGame(): boolean {
    return this.isPage(Page.Game);
  }
}