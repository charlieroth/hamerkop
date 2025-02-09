import app from './src/app.ts';
import { logi } from '@soapbox/logi';
import Config from './src/config.ts';

Deno.serve({
  port: Config.port,
  onListen({ hostname, port }) {
    logi({
      level: 'info',
      ns: 'hamerkop.server',
      message: `🦤 Hamerkop is running on ${hostname}:${port}`,
    });
  },
}, app.fetch);
