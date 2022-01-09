# caddy_dnsimple

## Usage

Using [xcaddy](https://github.com/caddyserver/xcaddy), you can build Caddy with this plugin with `xcaddy build --with github.com/hack-as-a-service/caddy_dnsimple`.

This plugin doesn't support `Caddyfile` usage, though you can use it in your JSON configuration:

```jsonc
{
  "apps": {
    "tls": {
      "automation": {
        "policies": [
          {
            "subjects": ["*.hackclub.app"],
            "issuers": [
              {
                "module": "acme",
                "challenges": {
                  "dns": {
                    "provider": {
                      "name": "dnsimple",
                      "api_token": "YOUR_API_TOKEN",
                      "account_id": "YOUR_ACCOUNT_ID"
                    }
                  }
                }
              }
            ]
          }
        ]
      }
    },
    "http": {
      "servers": {
        "srv0": {
          // ...
        }
      }
    }
  }
}
```
