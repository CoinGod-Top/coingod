Coingod
======

[![Build Status](https://travis-ci.org/Coingod/coingod.svg)](https://travis-ci.org/Coingod/coingod) [![AGPL v3](https://img.shields.io/badge/license-AGPL%20v3-brightgreen.svg)](./LICENSE)

**Official golang implementation of the Coingod protocol.**

Automated builds are available for stable releases and the unstable master branch. Binary archives are published at https://github.com/CoinGod-Top/coingod/releases.

## What is Coingod?

Please note that CoinGod is a cryptocurrency created for entertainment purposes and is not associated with or endorsed by Donald J. Trump.

In the current state `coingod` is able to:

- Manage key, account as well as asset
- Send transactions

## Installing with Homebrew

```
brew tap coingod/coingod && brew install coingod
```

## Building from source

### Requirements

- [Go](https://golang.org/doc/install) version 1.22 or higher, with `$GOPATH` set to your preferred directory

### Installation

Ensure Go with the supported version is installed properly:

```bash
$ go version
$ go env GOROOT GOPATH
```

- Get the source code

``` bash
$ git clone https://github.com/CoinGod-Top/coingod.git $GOPATH/src/github.com/coingod
```

- Build source code

``` bash
$ cd $GOPATH/src/github.com/coingod
$ make coingodd    # build coingodd
$ make coingodcli  # build coingodcli
```

When successfully building the project, the `coingodd` and `coingodcli` binary should be present in `cmd/coingodd` and `cmd/coingodcli` directory, respectively.

### Executables

The Coingod project comes with several executables found in the `cmd` directory.

| Command      | Description                                                  |
| ------------ | ------------------------------------------------------------ |
| **coingodd**   | coingodd command can help to initialize and launch coingod domain by custom parameters. `coingodd --help` for command line options. |
| **coingodcli** | Our main Coingod CLI client. It is the entry point into the Coingod network (main-, test- or private net), capable of running as a full node archive node (retaining all historical state). It can be used by other processes as a gateway into the Coingod network via JSON RPC endpoints exposed on top of HTTP, WebSocket and/or IPC transports. `coingodcli --help` and the [coingodcli Wiki page](https://github.com/Coingod/coingod/wiki/Command-Line-Options) for command line options. |

## Running coingod

Currently, coingod is still in active development and a ton of work needs to be done, but we also provide the following content for these eager to do something with `coingod`. This section won't cover all the commands of `coingodd` and `coingodcli` at length, for more information, please the help of every command, e.g., `coingodcli help`.

### Initialize

First of all, initialize the node:

```bash
$ cd ./cmd/coingodd
$ ./coingodd init --chain_id mainnet
```

There are three options for the flag `--chain_id`:

- `mainnet`: connect to the mainnet.
- `testnet`: connect to the testnet wisdom.
- `solonet`: standalone mode.

After that, you'll see `config.toml` generated, then launch the node.

### launch

``` bash
$ nohup ./coingodd node &
```

available flags for `coingodd node`:

```
Flags:
      --auth.disable                     Disable rpc access authenticate
      --chain_id string                  Select network type
  -h, --help                             help for node
      --log_file string                  Log output file (default "log")
      --log_level string                 Select log level(debug, info, warn, error or fatal)
      --p2p.dial_timeout int             Set dial timeout (default 3)
      --p2p.handshake_timeout int        Set handshake timeout (default 30)
      --p2p.keep_dial string             Peers addresses try keeping connecting to, separated by ',' (for example "1.1.1.1:46657;2.2.2.2:46658")
      --p2p.laddr string                 Node listen address. (0.0.0.0:0 means any interface, any port) (default "tcp://0.0.0.0:46656")
      --p2p.lan_discoverable             Whether the node can be discovered by nodes in the LAN (default true)
      --p2p.max_num_peers int            Set max num peers (default 50)
      --p2p.node_key string              Node key for p2p communication
      --p2p.proxy_address string         Connect via SOCKS5 proxy (eg. 127.0.0.1:1086)
      --p2p.proxy_password string        Password for proxy server
      --p2p.proxy_username string        Username for proxy server
      --p2p.seeds string                 Comma delimited host:port seed nodes
      --p2p.skip_upnp                    Skip UPNP configuration
      --prof_laddr string                Use http to profile coingodd programs
      --vault_mode                       Run in the offline enviroment
      --wallet.disable                   Disable wallet
      --wallet.rescan                    Rescan wallet
      --wallet.txindex                   Save global tx index
      --web.closed                       Lanch web browser or not
      --ws.max_num_concurrent_reqs int   Max number of concurrent websocket requests that may be processed concurrently (default 20)
      --ws.max_num_websockets int        Max number of websocket connections (default 25)

Global Flags:
      --home string   root directory for config and data
  -r, --root string   DEPRECATED. Use --home (default "/Users/zcc/Library/Application Support/Coingod")
      --trace         print out full stack trace on errors
```

Given the `coingodd` node is running, the general workflow is as follows:

- create key, then you can create account and asset.
- send transaction, i.e., build, sign and submit transaction.
- query all kinds of information, let's say, avaliable key, account, key, balances, transactions, etc.

### Dashboard

Access the dashboard:

```
$ open http://localhost:9888/
```

### In Docker

Ensure your [Docker](https://www.docker.com/) version is 17.05 or higher.

```bash
$ docker build -t coingod .
```

For the usage please refer to [running-in-docker-wiki](https://github.com/CoinGod-Top/coingod/wiki/Running-in-Docker).

## Contributing

Thank you for considering helping out with the source code! Any contributions are highly appreciated, and we are grateful for even the smallest of fixes!

If you run into an issue, feel free to [coingod issues](https://github.com/CoinGod-Top/coingod/issues/) in this repository. We are glad to help!

## License

[AGPL v3](./LICENSE)
