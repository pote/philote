# Philote:  plug-and-play websockets server

Philote is a minimal solution to the websockets server problem

### Notice:

Philote used to be powered by Redis's PUB/SUB capabilities, it's been rewritten to do all the work itself (it takes, surprisingly, less code to do the latter), this project is still in alpha so please handle with care.


## Deploy your own instance

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

### Installing it

A homebrew package is in the works, for now you can check out the [latest release](https://github.com/pote/philote/releases) and download the appropriate binary, or [install from source](#install-from-source)

### Running the server

There are three configuration options for Philote, all of which have sensible defaults and can be set by setting their respective environment variable.

| Environment Variable    | Default                   | Description                                    |
|:-----------------------:|:-------------------------:|:-----------------------------------------------|
| `PORT`                  | `6380`                    | Port in which to serve websocket connections   | 

If the defaults work for you, simply running `philote` will start the server with the default values, or you can just manipulate the environment and run with whatever settings you need.

```bash
$ PORT=9424 philote
```

## Clients

The only official client for Philote at the time of this writing is [philote-js](https://github.com/13floor/philote-js), as the primary use case are browser clients. Philote websocket's facilities are fairly standard though, so I expect clients in other languages should be straightforward to implement.

## Authentication

## Local development

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

## License

Released under MIT License, check LICENSE file for details.
