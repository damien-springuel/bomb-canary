import Axios from 'axios';
import App from './App.svelte';
import { MessageBus } from './messages/messagebus';
import { Party } from './party/party';
import { Store } from './store/store';
import { Handler } from './stream/handler';
import { Opener } from './stream/opener';
import { Creator } from './stream/creator';
import { PageManager } from './consumers/page';
import { ReplayManager } from './consumers/replay';
import { PlayerManager } from './consumers/player';
import { PlayerActions } from './player-actions/player-actions';
import { ResetManager } from './consumers/reset';

const hostname = window.location.hostname;

const axiosInstance = Axios.create({baseURL: `http://${hostname}:44324`, withCredentials: true});

const store = new Store();

const messageBus = new MessageBus();
messageBus.subscribeConsumer({
  consume: (m) => console.log(`Incoming Message: `, m)
});

const party = new Party(axiosInstance, messageBus);
messageBus.subscribeConsumer(party);

const playerActions = new PlayerActions(axiosInstance);
messageBus.subscribeConsumer(playerActions);

const handler = new Handler(messageBus);
const creator = new Creator(() => new WebSocket(`ws://${hostname}:44324/events`), handler);
const opener = new Opener(creator);
messageBus.subscribeConsumer(opener);

messageBus.subscribeConsumer(new ResetManager(store));
messageBus.subscribeConsumer(new PageManager(store));
messageBus.subscribeConsumer(new ReplayManager(store));
messageBus.subscribeConsumer(new PlayerManager(store));

const app = new App({
  target: document.body,
  props: {
    dispatcher: messageBus,
    store: store,
  }
});

export default app;