import type { AppController } from '../app.ts';
import Relay from '../relay.ts';

export const relayController: AppController = async (c, next) => {
  const ip = c.req.header('x-real-ip');
  const { socket, response } = Deno.upgradeWebSocket(c.req.raw, { idleTimeout: 30 });
  const kv = await Deno.openKv();
  const relay = new Relay(kv);
  relay.handleSocket(socket, ip);
  return response;
};
