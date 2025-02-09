import { Context, Env as HonoEnv, Handler, Hono, Input as HonoInput, MiddlewareHandler } from '@hono/hono';
import { relayInfoController } from './controllers/relay-info-controller.ts';
import { relayController } from './controllers/relay-controller.ts';

export interface AppEnv extends HonoEnv {}
export type AppContext = Context<AppEnv>;
export type AppMiddleware = MiddlewareHandler<AppEnv>;
export type AppController<P extends string = any> = Handler<AppEnv, P, HonoInput, Response | Promise<Response>>;

const app = new Hono<AppEnv>({ strict: false });

app.get('/', (c, next) => {
  const upgrade = c.req.header('upgrade');

  // NIP-11
  if (c.req.header('accept') === 'application/nostr+json') {
    return relayInfoController(c, next);
  }

  if (upgrade?.toLowerCase() !== 'websocket') {
    return c.text('Please use a Nostr client to connect to this relay.', 400);
  }

  return relayController(c, next);
});

export default app;
