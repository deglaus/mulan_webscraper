// Exports subscriber the instance of BehaviorSubject and the messageService object

import { BehaviorSubject } from "rxjs";

const subscriber = new BehaviorSubject(0);

const messageService = {
  send: function(msg) {
    subscriber.next(msg);
  },
};

export { messageService, subscriber };
