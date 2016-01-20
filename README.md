# Philote: A Redis-powered websockets server.

Philote is a minimal solution to the websockets server problem, it doesn't even do most of the work: it acts as a bridge between websockets clients such as browser JavaScript engines and a [Redis](http://redis.io/) instance, taking advantage of it's [PubSub](http://redis.io/commands#pubsub) capabilities.

Philote has almost zero-configuration, as it already relies on Redis, websocket clients identify themselves with a token, you can create this tokens in your applications and store them in Redis, which will determine the level of access that connection will have to different pubsub channels - more on this later.

## Bootstrap it

You'll need [Redis](http://redis.io/) installed and running, and [gpm](https://github.com/pote/gpm) for dependency management.

``` bash
$ source .env.sample # you might want to copy it to .env if you plan on changing the settings)
$ make
```

### Run it.

``` bash
$ make server
```

### Run the test suite

```bash
$ make test
```

### `philote-admin`

You'll find an executable in `admin/philote-admin`, it's mainly a development help, it can create access keys and publish messages to channels so you can try philote locally with ease.

Run `./admin/philote-admin --help` for more.

## Clients

The only official client for Philote at the time of this writing is [philote-js](https://github.com/13floor/philote-js), as the primary use case are browser clients. Philote websocket's facilities are fairly standard though, so I expect clients in other languages should be straightforward to implement.

## Authentication

Given a `$PHILOTE_URL` of where your Philote server lives, clients will need to connect to `$PHILOTE_URL/<identifier-token>`, Philote uses the `identifier-token` to look for an access key in Redis (namespace `philote:access_key:<identifier-token>`), if the token is valid whatever permissions specified in the Redis keys will be applied to the connection.

## Access Keys

Access keys are stored in Redis (under the key `philote:access_key:<identifier-token>`), and consist of a JSON-encoded Hash, which looks like this:

```json
{
  "read": [
    "test-channel"
  ],
  "write": [
    "test-channel"
  ],
  "allowed_uses": 1,
  "uses": 0
}
``` 
You can create your own access keys in your language of choice and store them in Redis with whatever identifier you choose, ranging from secure, randomly-generated tokens to simpler, straightforward ones, keep in mind that these tokens will be exposed to the public so exercise as much caution as your use case requires, for most cases we recommend making each token a single-use one, setting `allowed_uses` to `1` as per the example, you can also [specify an expire time](http://redis.io/commands/set) when you store it in Redis. 

We'll make tools that make it easy to create tokens for common languages such as Ruby and Python, or other languages we use regularly like Lua and Go in due time, but really you can use Philote very easily by generating and storing the access keys in your language of choice.

# Access Key JSON Schema

There's a JSON schema for the AccessKeys [included in this repo](./meta/access-key-schema.json).

```json
{
  "$schema": "http://json-schema.org/draft-04/schema#",
  "title": "Philote's Access Key",
  "description": "Describes permissions for Philote websocket connection",
  "type": "object",
  "properties": {
    "read": {
      "description": "PubSub channels for which the connection will receive messages",
      "type": "array",
      "items": {
        "type": "string"
      },
      "uniqueItems": true
    },
    "write": {
      "description": "PubSub channels for which the connection will be allowed to publish messages",
      "type": "array",
      "items": {
        "type": "string"
      },
      "uniqueItems": true
    },
    "allowed_uses": {
      "description": "Amount of times this access key can be used to connect to Philote (0 means unlimited usage)",
      "type": "integer"
    },
    "uses": {
      "description": "Amount of times this access key was used to connect to Philote",
      "type": "integer"
    }
  },
  "required": ["read", "write", "allowed_uses", "uses"]
}
```
## Caveats

Philote opens a connection to Redis per websocket connection that it maintains open, if open connections are a limitation of your main Redis database I'd recommend having Philote use a separate one.

## License

Released under MIT License, check LICENSE file for details.

## Sponsorship

This open source tool is proudly sponsored by [13Floor](http://13Floor.org)

![13Floor](./13Floor-circulo-1.png)