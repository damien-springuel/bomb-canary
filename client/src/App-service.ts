import type { Dispatcher } from "./messages/dispatcher";
import { AppLoaded } from "./messages/events";

export class AppService {
  constructor(
    private readonly dispatcher: Dispatcher){}

  appMounted() {
    this.dispatcher.dispatch(new AppLoaded());
  }
}