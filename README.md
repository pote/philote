# Philote:  plug-and-play websockets server ![Build status](https://travis-ci.org/pote/philote.svg)

Philote is a minimal solution to the websockets server problem, it implements Publish/Subscribe and has a simple authentication mechanism that accomodates browser clients securely as well as server-side or even local applications.

Simplicity is one of the design goals for Philote, ease of deployment is another: you should be able to drop the binary in any internet-accessible server and have it operational.

For a short demonstration, check out the sample command line Philote client called [Jane](#cli)

## Basics

Philote implements a basic topic-based [Publish-subscribe pattern](https://en.wikipedia.org/wiki/Publish%E2%80%93subscribe_pattern), messages sent over the websocket connection are classified into `channels`, and each connection is given read/write access to a given list of channels at authentication time.

Messages sent over a connection for a given channel (to which it has write permission) will be received by all other connections (that have read permission to the channel in question).

### Deploy your own instance

You can play around with Philote by deploying it on Heroku for free, keep in mind that Heroku's free tier dynos are not suited for production Philote usage, however, as sleeping dynos will mean websocket connections are closed.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

### Configuration options

There are three configuration options for Philote, all of which have sensible defaults and can be set by setting their respective environment variable.

| Environment Variable    | Default                   | Description                                                                                                        |
|:-----------------------:|:-------------------------:|:-------------------------------------------------------------------------------------------------------------------|
| `SECRET`                | ` `                       | Secret salt used to sign authentication tokens                                                                     |
| `PORT`                  | `6380`                    | Port in which to serve websocket connections                                                                       |
| `LOGLEVEL`              | `info`                    | Verbosity of log output, valid options are [debug,info,warning,error,fatal,panic]                                  |
| `MAX_CONNECTIONS`       | `255`                     | Maximum amount of concurrent websocket connections allowed                                                         |
| `READ_BUFFER_SIZE`      | `1024`                    | Size of the websocket read buffer, for most cases the default should be okay.                                      |
| `WRITE_BUFFER_SIZE`     | `1024`                    | Size of the websocket write buffer, for most cases the default should be okay.                                     |
| `CHECK_ORIGIN`          | `false`                   | Check Origin headers during WebSocket upgrade handshake.                                                           |

If the defaults work for you, simply running `philote` will start the server with the default values, or you can just manipulate the environment and run with whatever settings you need.

```bash
$ PORT=9424 philote
```

## CLI

There is a trivial implementation of basic Philote interaction called [Jane](https://github.com/pote/jane) that you can run locally, it can subscribe to a channel on a Philote server, receive and publish messages. It's useful for debugging purposes.

![sample](https://stuff.pote.io/Screen-Recording-2017-05-16-15-50-30-5ivJp0cbze.gif)

## Clients

* [JavaScript (browser)](https://github.com/pote/philote-js)
* [Go](https://github.com/pote/philote-go)

## Authentication

Clients authenticate in Philote using [JSON Web Tokens](https://jwt.io), which consist on a JSON payload detailing the read/write permissions a given connection will have. The payload is hashed with a secret known to Philote so that incoming connections can be verified, this way you can generate tokens in your application backend and use them from the browser client without fear.

Clients in different language will provide methods to generate these tokens, for now, the [Go client](https://github.com/pote/philote-go/blob/master/token.go) should be the reference implementation, although you'll notice that it's an extremely simple one so ports to other languages should be trivial to implement provided with a decent JWT library.

For incoming websockets connections, Philote will look to find the authentication token in the `Authorization` header, but since the native browser JavaScript WebSocket API does not provide a way to manipulate the request headers Philote will also look for the `auth` query parameter in case it fails to authenticate using the header option.


### Install

You can install Philote (and Jane) easily with homebrew.

`brew install pote/philote/philote`

`brew install pote/philote/jane`

You can also manually get the binaries from [latest release](https://github.com/pote/philote/releases) or [install from source](#install-from-source)


### Local development

### Bootstrap it

You'll need [gpm](https://github.com/pote/gpm) for dependency management.

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

## History

The first versions of Philote were powered by Redis, it was initially thought of as a websocket bridge to a Redis instance.

After a while, that design was considered inpractical: redis is a big dependency to have, publish/subscribe was easy to implement in Philote itself and the authentication mechanism was changed to use JSON Web Tokens, making Redis unnecessary.

The result should be a more robust tool that anyone can drop in any operating system and get working trivially, without external dependencies.

## License

Released under MIT License, check LICENSE file for details.
