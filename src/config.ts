export default class Config {
  static get port(): number {
    return 8000
  }

  static get host(): string {
    return "127.0.0.1"
  }

  static get localDomain(): string {
    return Deno.env.get('LOCAL_DOMAIN') || `http://localhost:${Config.port}`
  }

  static get url(): URL {
    return new URL(Config.localDomain)
  }

  static get relay(): string {
    const { protocol, host } = Config.url;
    return `${protocol === 'https:' ? 'wss:' : 'ws:'}//${host}/relay`;
  }
}
