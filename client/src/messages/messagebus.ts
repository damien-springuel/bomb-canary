export interface Message {}

export interface Event extends Message {}

export interface Consumer{
  consume(message: Message): void
}

export class MessageBus {
  private readonly consumers: Consumer[] = [];

  public subscribeConsumer(consumer: Consumer) {
    this.consumers.push(consumer);
  }

  public dispatch(message: Message): void {
    this.consumers.forEach(c => c.consume(message));
  }
}
