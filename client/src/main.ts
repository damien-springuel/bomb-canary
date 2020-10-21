import Axios from 'axios';
import App from './App.svelte';
import { MessageBus } from './messages/messagebus';
import { Party } from './party/party';
import { Store } from './store/store';
import { Handler } from './stream/handler';
import { ServerEventConnectionOpener } from './stream/server-event-connection';
import { ServerStream } from './stream/server-stream';

const axiosInstance = Axios.create({baseURL: "http://localhost:44324", withCredentials: true});

const store = new Store();

const messageBus = new MessageBus();
messageBus.subscribeConsumer({
  consume: (m) => console.log(`Incoming Message: `, m)
});

const party = new Party(axiosInstance, messageBus);
messageBus.subscribeConsumer(party);

const handler = new Handler(messageBus);
const serverStream = new ServerStream(() => new WebSocket(`ws://localhost:44324/events`), handler);
const serverEventConnectionOpener = new ServerEventConnectionOpener(serverStream);
messageBus.subscribeConsumer(serverEventConnectionOpener);

const app = new App({
  target: document.body,
  props: {
    dispatcher: messageBus,
    store: store,
  }
});

export default app;