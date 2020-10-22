import test from "ava";
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

test(`Creator - on close`, t => {
  let {websocket, handler} = setup();
  websocket.onclose({} as CloseEvent);
  t.true(handler.onCloseCalled);
});

test(`Creator - on error`, t => {
  let {websocket, handler} = setup();
  websocket.onerror({} as Event);
  t.true(handler.onErrorCalled);
});

test(`Creator - on message`, t => {
  let {websocket, handler} = setup();
  const testServerEvent = {test: `event`} as ServerEvent
  websocket.onmessage({data: JSON.stringify(testServerEvent)} as MessageEvent);
  t.deepEqual(handler.receivedEvent, testServerEvent);
});
