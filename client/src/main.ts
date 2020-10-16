import App from './App.svelte';
import { Message, MessageBus } from './messages/messagebus';
import { Store } from './store/store';

const store = new Store();

const messageBus = new MessageBus();
messageBus.SubscribeConsumer({
  consume: (m) => console.log(`Incoming Message: `, m)
})

const app = new App({
  target: document.body,
  props: {
    dispatcher: messageBus,
    store: store,
  }
});

export default app;