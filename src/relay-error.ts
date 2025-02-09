import { NostrRelayOK } from '@nostrify/nostrify';

export type RelayErrorPrefix = 'duplicate' | 'pow' | 'blocked' | 'rate-limited' | 'invalid' | 'error';

// NIP-01 result
export class RelayError extends Error {
  constructor(prefix: RelayErrorPrefix, message: string) {
    super(`${prefix}: ${message}`);
  }

  static fromReason(reason: string): RelayError {
    const [prefix, ...rest] = reason.split(': ');
    return new RelayError(prefix as RelayErrorPrefix, rest.join(': '));
  }

  static assert(msg: NostrRelayOK): void {
    const [, , ok, reason] = msg;
    if (!ok) {
      throw RelayError.fromReason(reason);
    }
  }
}
