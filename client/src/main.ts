import Axios from 'axios';
import App from './App.svelte';
import { Message, MessageBus } from './messages/messagebus';
import { Party } from './party/party';
import { Store } from './store/store';

const axiosInstance = Axios.create({baseURL: "http://localhost:44324", withCredentials: true});

const store = new Store();

const messageBus = new MessageBus();
messageBus.subscribeConsumer({
  consume: (m) => console.log(`Incoming Message: `, m)
});

const party = new Party(axiosInstance, messageBus);
messageBus.subscribeConsumer(party);

const app = new App({
  target: document.body,
  props: {
    dispatcher: messageBus,
    store: store,
  }
});

export default app;