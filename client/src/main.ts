import App from './App.svelte';
import { CreatePartyClicked, Message, MessageBus } from './messagebus';
import { Store } from './store/store';

const store = new Store();

const messageBus = new MessageBus();
messageBus.SubscribeConsumer({
  consume: (m) => console.log(`Incoming Message: `, m)
})


class Party {
  constructor(private readonly nameUpdater: {setName:(name:string) => void}) {}

  consume(message: Message) {
    if(message instanceof CreatePartyClicked) {
      this.nameUpdater.setName('Cool ' + message.name);
    }
  }
}

const party = new Party(store);

messageBus.SubscribeConsumer(party);

const app = new App({
  target: document.body,
  props: {
    dispatcher: messageBus,
    store: store,
  }
});

export default app;