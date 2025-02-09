// Serialize an error to a JSON object
export function errorJson(error: unknown): Error | null {
  if (error instanceof Error) {
    return error;
  } else {
    return null;
  }
}
