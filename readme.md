# Traefik XFF Fix Plugin

A [Traefik](https://traefik.io) middleware plugin that rewrites the `X-Forwarded-For` header to contain only the leftmost (original client) IP address, with any port number stripped. This is useful when downstream services expect a single clean client IP instead of a comma-separated chain.

## Behavior

| Input | Output |
|-------|--------|
| `203.0.113.195, 70.41.3.18, 150.172.238.178` | `203.0.113.195` |
| `203.0.113.195:8080, 70.41.3.18` | `203.0.113.195` |
| `[2001:db8::1]:8080, 10.0.0.1` | `2001:db8::1` |
| `2001:db8::1, 10.0.0.1` | `2001:db8::1` |
| `192.168.1.1` | `192.168.1.1` |

## Configuration

### Static configuration

Declare the plugin in your Traefik static configuration:

```yaml
experimental:
  plugins:
    xff-fix:
      moduleName: github.com/marcelfzr/traefik-xff-fix
      version: v0.1.0
```

### Dynamic configuration

Attach the middleware to your routers:

```yaml
http:
  routers:
    my-router:
      rule: host(`example.com`)
      service: my-service
      middlewares:
        - xff-fix

  middlewares:
    xff-fix:
      plugin:
        xff-fix: {}
```

### Local development

For testing plugins locally without publishing to GitHub:

```yaml
experimental:
  localPlugins:
    xff-fix:
      moduleName: github.com/marcelfzr/traefik-xff-fix
```

Place the plugin source at:

```
./plugins-local/src/github.com/marcelfzr/traefik-xff-fix/
```

## Plugins Catalog

This plugin is hosted at [github.com/marcelfzr/traefik-xff-fix](https://github.com/marcelfzr/traefik-xff-fix). Once the repository has the `traefik-plugin` topic and a valid `.traefik.yml` manifest, it will be discoverable in the [Traefik Plugins Catalog](https://plugins.traefik.io).
