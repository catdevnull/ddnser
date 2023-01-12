# ddnser

Refreshes domain records periodically or when interfaces change IP addresses.

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
    }
  ]
}
```

Run:

```sh
ddnser ./path/to/config.json
```
