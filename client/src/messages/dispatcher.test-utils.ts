import type { Message } from "./messagebus";

export class AsyncDispatcherMock {
  public receivedMessage: Message
  public readonly isDone: Promise<void>
  private done: (value?: void) => void

  constructor(){
    this.isDone = new Promise((resolver, rejecter) => {
      this.done = resolver;
    });
  }

  dispatch(m: Message): void {
    this.receivedMessage = m;
    this.done();
  }
}