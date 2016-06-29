# Philote: A Redis-powered websockets server.

Philote is a minimal solution to the websockets server problem, it doesn't even do most of the work: it acts as a bridge between websockets clients such as browser JavaScript engines and a [Redis](http://redis.io/) instance, taking advantage of it's [PubSub](http://redis.io/commands#pubsub) capabilities.

## What it does

Philote has two features: it serves websockets connections and it provides an authentication mechanism for clients opening them.

Authentication happens through [Access Keys](#access-keys), these keys are stored in Redis and identified by a token which the client will need in order to open a connection to Philote.

Once open, a websocket connection has read and/or write access to a set of channels, they will receive messages pushed by other clients into those channels and have the ability to publish their own messages.

Philote pub/sub capabilities are backed by redis's own, what this means is you can interact with the philote channels simply by publishing or listening to pub/sub messages in redis, without a need for special clients and without the overhead of opening websockets connections, publishing is as fast as sending out a Redis [PUBLISH](http://redis.io/commands/publish) command.

### Installing it

A homebrew package is in the works, for now you can check out the [latest release](https://github.com/pote/philote/releases) and download the appropriate binary, or [install from source](#install-from-source)

### Running the server

There are three configuration options for Philote, all of which have sensible defaults and can be set by setting their respective environment variable.

| Environment Variable    | Default                   | Description                                    |
|:-----------------------:|:-------------------------:|:-----------------------------------------------|
| `PORT`                  | `6380`                    | Port in which to serve websocket connections   | 
| `REDIS_URL`             | `redis://localhost:6379`  | Philote-backing Redis instance                 | 
| `REDIS_MAX_CONNECTIONS` | `400`                     | Maximum number of concurrent Redis connections |

If the defaults work for you, simply running `philote` will start the server with the default values, or you can just manipulate the environment and run with whatever settings you need.

```bash
$ PORT=9424 REDIS_URL=redis://123.412.512.3:12352/2 philote
```

### philote-admin

You'll find an executable called `philote-admin`, it's mainly a development help, it can create access keys and publish messages to channels so you can try philote locally with ease.

Run `philote-admin --help` for more.


## Clients

The only official client for Philote at the time of this writing is [philote-js](https://github.com/13floor/philote-js), as the primary use case are browser clients. Philote websocket's facilities are fairly standard though, so I expect clients in other languages should be straightforward to implement.

There is also a [Ruby](https://github.com/pote/philote-rb) client, which currently only implements admin capabilities (creating access keys and publishing via Redis).

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

If you're using Ruby, we've made [a helper library](https://github.com/pote/philote-rb) that should make creating access keys and publishing messages simpler for you. We'll make tools that make it easy to create tokens for common languages such as Python, Lua, and Go in due time, but really you can use Philote very easily by generating and storing the access keys in your language of choice and storing it directly in Redis.

There's a JSON schema for the AccessKeys [included in this repo](./meta/access-key-schema.json).

## Caveat

Philote maintains an open Redis connection for each websocket connection that it serves, so it is recommended to use a dedicated Redis instance to back Philote.

This is likely going to be refactored into something smarter in the near future, but before doing that some stress stesting is needed to figure out the limits of Redis concurrent connection support as well as observing real world usage to find out how much of a problem it is.

## Local development

### Bootstrap it

You'll need [Redis](http://redis.io/) installed and running, and [gpm](https://github.com/pote/gpm) for dependency management.

### Set it up

``` bash
$ source .env.sample # you might want to copy it to .env and source that instead if you plan on changing the settings.
$ make 
```

### Run a Philote server

```bash
$ make server
```

### Run the test suite

```bash
$ make test
```

### Install from source

```bash
$ make install
```

## Further reading

There is an [example chat application](https://github.com/pote/philote-chat-app), powered by Philote, [philote-rb](https://github.com/pote/philote-rb) and [philote-js](https://github.com/13Floor/philote-js) which you can look at for a peek into how you could implement Philote into your app.

## License

Released under MIT License, check LICENSE file for details.
