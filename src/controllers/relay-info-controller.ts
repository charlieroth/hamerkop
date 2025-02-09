import type { AppController } from '../app.ts';
import denoJson from 'deno.json' with { type: 'json' };

export const relayInfoController: AppController = (c) => {
  return c.json({
    name: 'hamerkop',
    description: 'A Nostr relay for client specializing in long-form content.',
    version: denoJson.version,
    contact: 'https://charlieroth.me',
    software: 'https://github.com/charlieroth/hamerkop',
    supported_nips: [
      1,
      // 2,
      // 9,
      // 11,
      // 22,
      // 23,
      // 25,
      // 30,
      // 36,
      // 40,
      // 42,
      // 45,
      // 47,
      // 50,
      // 51,
      // 56,
      // 57,
      // 58,
      // 60,
      // 61,
    ],
    limitation: {
      auth_required: false,
      created_at_lower_limit: 0,
      created_at_upper_limit: 2_147_483_647,
      max_limit: 100,
      payment_required: false,
      restricted_writes: false,
    },
  });
};
