import type { Message } from "./messagebus";

export interface Dispatcher {
  dispatch(message: Message): void
}