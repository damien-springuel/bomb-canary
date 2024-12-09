import Axios from 'axios';
import App from './App.svelte';
import { MessageBus } from './messages/message-bus';
import { Party } from './party/party';
import { Store } from './store/store';
import { Handler } from './stream/handler';
import { Opener } from './stream/opener';
import { Creator } from './stream/creator';
import { PageConsumer } from './consumers/page';
import { ReplayConsumer } from './consumers/replay';
import { PlayerConsumer } from './consumers/player';
import { PlayerActions } from './player-actions/player-actions';
import { ResetConsumer } from './consumers/reset';
import { GameConsumer } from './consumers/game';

const axiosInstance = Axios.create({baseURL: window.location.origin, withCredentials: true});

const store = new Store();

const messageBus = new MessageBus();
if(import.meta.env.DEV) {
  messageBus.subscribeConsumer({
    consume: (m) => console.log(`Incoming Message: `, m)
  });
}

const party = new Party(axiosInstance, messageBus);
messageBus.subscribeConsumer(party);

const playerActions = new PlayerActions(axiosInstance);
messageBus.subscribeConsumer(playerActions);

const handler = new Handler(messageBus);
const creator = new Creator(() => new WebSocket(`${window.location.origin.replace("http", "ws")}/events`), handler);
const opener = new Opener(creator);
messageBus.subscribeConsumer(opener);

messageBus.subscribeConsumer(new ResetConsumer(store));
messageBus.subscribeConsumer(new PageConsumer(store));
messageBus.subscribeConsumer(new ReplayConsumer(store));
messageBus.subscribeConsumer(new PlayerConsumer(store));
messageBus.subscribeConsumer(new GameConsumer(store));

const app = new App({
  target: document.body,
  props: {
    dispatcher: messageBus,
    store: store,
  }
});

export default app;