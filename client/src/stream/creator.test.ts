import { expect, test } from "vitest";
import type { ServerEvent } from "./server-event";
import { Creator } from "./creator";


class WebsocketMock{
  public onclose: (ev: CloseEvent) => any;
  public onerror: (ev: Event) => any;
  public onmessage: (ev: MessageEvent) => any;
}

class EventHandlerMock {
  public onCloseCalled: boolean;
  public onErrorCalled: boolean;
  public receivedEvent: ServerEvent

  onClose(): void {
    this.onCloseCalled = true;
  }
  onError(): void {
    this.onErrorCalled = true
  }
  onEvent(event: ServerEvent) {
    this.receivedEvent = event;
  }
}

function setup(): {websocket: WebsocketMock, handler: EventHandlerMock} {
  let websocket = new WebsocketMock();
  let handler = new EventHandlerMock();
  const ss = new Creator(() => websocket, handler);
  ss.create();
  return {websocket, handler};
}

test(`Creator - on close`, () => {
  let {websocket, handler} = setup();
  websocket.onclose({} as CloseEvent);
  expect(handler.onCloseCalled).to.be.true;
});

test(`Creator - on error`, () => {
  let {websocket, handler} = setup();
  websocket.onerror({} as Event);
  expect(handler.onErrorCalled).to.be.true;
});

test(`Creator - on message`, () => {
  let {websocket, handler} = setup();
  const testServerEvent = {test: `event`} as ServerEvent
  websocket.onmessage({data: JSON.stringify(testServerEvent)} as MessageEvent);
  expect(handler.receivedEvent).to.deep.equal(testServerEvent)
});
