import { JsonValue } from '@std/json';
import { logi } from '@soapbox/logi';
import {
  NostrClientCLOSE,
  NostrClientCOUNT,
  NostrClientEVENT,
  NostrClientMsg,
  NostrClientREQ,
  NostrRelayMsg,
  NSchema as n,
} from '@nostrify/nostrify';
import { NDenoKv } from '@nostrify/denokv';
import { errorJson } from './utils.ts';
import { RelayError } from './relay-error.ts';
import Config from './config.ts';

export default class Relay {
  private filterLimit: number;
  private store: NDenoKv;
  private connections: Set<WebSocket>;
  private controllers: Map<string, AbortController>;

  constructor(kv: Deno.Kv) {
    this.filterLimit = 100;
    this.store = new NDenoKv(kv);
    this.connections = new Set();
    this.controllers = new Map();
  }

  handleSocket(socket: WebSocket, ip: string | undefined) {
    socket.onopen = () => {
      this.connections.add(socket);
    };

    socket.onmessage = (message: MessageEvent) => {
      if (typeof message.data !== 'string') {
        socket.close(1003, 'invalid message');
        return;
      }

      const result = n.json().pipe(n.clientMsg()).safeParse(message.data);
      if (result.success) {
        logi({
          level: 'trace',
          ns: 'ditto.relay.message',
          data: result.data as JsonValue,
        });
        this.handleClientMessage(socket, result.data);
      } else {
        this.send(socket, ['NOTICE', 'invalid message']);
      }
    };

    socket.onclose = () => {
      this.connections.delete(socket);
      for (const controller of this.controllers.values()) {
        controller.abort();
      }
    };
  }

  handleClientMessage(
    socket: WebSocket,
    clientMessage: NostrClientMsg,
  ) {
    switch (clientMessage[0]) {
      case 'REQ':
        this.handleRequest(socket, clientMessage);
        break;
      case 'EVENT':
        this.handleEvent(socket, clientMessage);
        break;
      case 'CLOSE':
        this.handleClose(socket, clientMessage);
        break;
      case 'COUNT':
        this.handleCount(socket, clientMessage);
        break;
    }
  }

  async handleRequest(
    socket: WebSocket,
    [_, subscriptionId, ...filters]: NostrClientREQ,
  ): Promise<void> {
    const controller = new AbortController();
    this.controllers.get(subscriptionId)?.abort();
    this.controllers.set(subscriptionId, controller);

    try {
      // TODO: Query event store
      for (const event of await this.store.query(filters, { limit: this.filterLimit, timeout: 1000 })) {
      }
      // TODO: Send the events to the client
    } catch (err) {
      // In the case of any error, we close the subscription
      if (err instanceof RelayError) {
        this.send(socket, ['CLOSED', subscriptionId, err.message]);
      } else if (err instanceof Error && err.message.includes('timeout')) {
        this.send(socket, ['CLOSED', subscriptionId, 'relay could not complete request in time']);
      } else {
        this.send(socket, ['CLOSED', subscriptionId, 'unknown error']);
      }
      // Remove the abort controller for this subscription
      this.controllers.delete(subscriptionId);
      return;
    }

    // Send the EOSE message to the client signaling that the request is complete
    this.send(socket, ['EOSE', subscriptionId]);
  }

  async handleEvent(
    socket: WebSocket,
    [_, event]: NostrClientEVENT,
  ): Promise<void> {
    const controller = new AbortController();
    const { signal } = controller;

    try {
      // TODO: Save the event to the event store through pipeline of filters
      // TODO: Send the event to the client
      // send(socket, ['OK', event.id, true, ''])
    } catch (err) {
      if (err instanceof RelayError) {
        this.send(socket, ['OK', event.id, false, err.message]);
      } else {
        this.send(socket, ['OK', event.id, false, 'unknown error']);
        logi({
          level: 'error',
          ns: 'hamerkop.relay',
          msg: 'Error in relay',
          error: errorJson(err),
        });
      }
    }
  }

  async handleClose(
    socket: WebSocket,
    [_, subscriptionId]: NostrClientCLOSE,
  ): Promise<void> {
    const controller = this.controllers.get(subscriptionId);
    if (controller) {
      controller.abort();
      this.controllers.delete(subscriptionId);
    }
  }

  async handleCount(
    socket: WebSocket,
    [_, subscriptionId, ...filters]: NostrClientCOUNT,
  ): Promise<void> {
    try {
      // TODO: Query event store
      // TODO: Send the count to the client
      // send(socket, ['COUNT', subscriptionId, { count, approximate: false }])
    } catch {
      this.send(socket, ['CLOSED', subscriptionId, 'unknown error']);
    }
  }

  send(socket: WebSocket, msg: NostrRelayMsg) {
    if (socket.readyState === WebSocket.OPEN) {
      socket.send(JSON.stringify(msg));
    }
  }
}
