import type { ServerEvent } from "./server-event";

interface BasicWebsocket {
  onclose: (ev: CloseEvent) => any;
  onerror: (ev: Event) => any;
  onmessage: (ev: MessageEvent) => any;
}

interface EventHandler {
  onClose(): void;
  onError(): void;
  onEvent(event: ServerEvent): void
}

export class Creator {
  constructor(
    private readonly wsCreator: () => BasicWebsocket,
    private readonly handler: EventHandler,
  ){}

  create(): void {
    let socket = this.wsCreator();

    socket.onmessage = event => {
      let gameEvent: ServerEvent = JSON.parse(event.data);
      this.handler.onEvent(gameEvent);
    };

    socket.onclose = () => {
      this.handler.onClose();
    };

    socket.onerror = () => {
      this.handler.onError();
    };
  }
}