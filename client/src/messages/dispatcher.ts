import type { Message } from "./message-bus";

export interface Dispatcher {
  dispatch(message: Message): void
}