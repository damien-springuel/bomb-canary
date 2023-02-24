import type { Message } from "./message-bus";

export class AsyncDispatcherMock {
  public receivedMessage: Message
  public readonly isDone: Promise<void>
  private done: (value?: void) => void

  constructor(){
    this.isDone = new Promise((resolver) => {
      this.done = resolver;
    });
  }

  dispatch(m: Message): void {
    this.receivedMessage = m;
    this.done();
  }
}

export class DispatcherMock {
  public receivedMessage: Message

  dispatch(m: Message): void {
    this.receivedMessage = m;
  }
}