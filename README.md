# Philote:  plug-and-play websockets server

Philote is a minimal solution to the websockets server problem, it implements Publish/Subscribe and has a simple authentication mechanism that accomodates browser clients securely as well as server-side or even local applications.

Simplicity is one of the design goals for Philote, ease of deployment is another: you should be able to drop the binary in any internet-accessible server and have it operational.


## Philote in action

This gif shows a [sample command line chat application](https://github.com/pote/jane) that connects to a Philote instance and enables users to exchange messages, this command line application uses the [Go philote client](#clients).

![sample](https://stuff.pote.io/Screen-Recording-2017-05-16-15-50-30-5ivJp0cbze.gif)

### Deploy your own instance

You can play around with Philote by deploying it on Heroku for free, keep in mind that Heroku's free tier dynos are not suited for production Philote usage, however, as sleeping dynos will mean websocket connections are closed.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

### Configuration options

There are three configuration options for Philote, all of which have sensible defaults and can be set by setting their respective environment variable.

| Environment Variable    | Default                   | Description                                                                                                        |
|:-----------------------:|:-------------------------:|:-------------------------------------------------------------------------------------------------------------------|
| `SECRET`                | ``                        | Secret salt used to sign authentication tokens, this secret needs to be known to the clients in order to connect   |
| `PORT`                  | `6380`                    | Port in which to serve websocket connections                                                                       |
| `LOGLEVEL`              | `info`                    | Verbosity of log output, valid options are [debug|info|warning|error|fatal|panic]                                  |
| `MAX_CONNECTIONS`       | `255`                     | Maximum amount of concurrent websocket connections allowed                                                         |
| `READ_BUFFER_SIZE`      | `1024`                    | Size of the websocket read buffer, for most cases the default should be okay.                                      |
| `WRITE_BUFFER_SIZE`     | `1024`                    | Size of the websocket write buffer, for most cases the default should be okay.                                      |

If the defaults work for you, simply running `philote` will start the server with the default values, or you can just manipulate the environment and run with whatever settings you need.

```bash
$ PORT=9424 philote
```

## Clients

## Authentication

## Local development

### Installing it

A homebrew package is in the works, for now you can check out the [latest release](https://github.com/pote/philote/releases) and download the appropriate binary, or [install from source](#install-from-source)

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
