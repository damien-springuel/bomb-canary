export interface Message {}

export interface Consumer{
  consume(message: Message): void
}

export class MessageBus {
  private readonly consumers: Consumer[] = [];

  public SubscribeConsumer(consumer: Consumer) {
    this.consumers.push(consumer);
  }

  public Dispatch(message: Message): void {
    this.consumers.forEach(c => c.consume(message));
  }
}

export class CreatePartyClicked implements Message {
  constructor(readonly name:string){}
}