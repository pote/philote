# Philote: Expose Redis PubSub via Websockets.

This app provides a thin layer on top of Redis PubSub so that you can generate
websocket connections that have two-way communication between them.

## Bootstrap it

``` bash
$ source .env.sample # you might want to copy it to .env if you plan on changing the settings)
$ gpm install
```

## Build it, run it.

``` bash
$ make
$ ./philote
```

The websocket server will then start running in `$PORT`, it expects to receive
connections at  `ws://localhost:$PORT/<hub_identifier>/<room>`, of course
connections that subscribe to the same room will receive each other's messages.


## See it in action

http://stuff.poteland.com/YvYDJ5HYha.mp4
