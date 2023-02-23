import { CloseDialog } from "../messages/commands";
import type { Dispatcher } from "../messages/dispatcher";

export class DialogService {
  constructor(private readonly dispatcher: Dispatcher){}

  closeDialog() {
    this.dispatcher.dispatch(new CloseDialog());
  }
}