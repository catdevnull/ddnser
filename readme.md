# ddnser

Refreshes domain records periodically.

## Usage

Create a config file:

```json
{
  // "ip": "hardcoded IP address, if set it is used instead of detected address",
  // how often to poll unconditionally
  "every": 60,
  "domains": [
    {
      "type": "njalla ddns",
      "name": "estoesprueba.nulo.in",
      "key": "INSERT_KEY"
    },
    {
      "type": "he.net ddns",
      "name": "pruebas.bat.ar",
      "key": "INSERT_KEY"
    },
    {
      "type": "cloudflare v4 api",
      "name": "*.nulo.in",
      "zoneName": "nulo.in",
      "key": "INSERT_KEY" // https://dash.cloudflare.com/profile/api-tokens
    }
  ]
}
```

Run:

```sh
ddnser ./path/to/config.json
```

## Docker

Compose:

```
  ddnser:
    image: ghcr.io/catdevnull/ddnser
    restart: always
    volumes:
      - ./ddnser.json:/config/ddnser.json:ro
    entrypoint: ddnser /config/ddnser.json
```
