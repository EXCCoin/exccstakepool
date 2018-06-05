# exccstakepool

[![GoDoc](https://godoc.org/github.com/EXCCoin/exccstakepool?status.svg)](https://godoc.org/github.com/EXCCoin/exccstakepool)
[![Build Status](https://travis-ci.org/EXCCoin/exccstakepool.svg?branch=master)](https://travis-ci.org/EXCCoin/exccstakepool)
[![Go Report Card](https://goreportcard.com/badge/github.com/EXCCoin/exccstakepool)](https://goreportcard.com/report/github.com/EXCCoin/exccstakepool)

exccstakepool is a web application which coordinates generating 1-of-2 multisig
addresses on a pool of [exccwallet](https://github.com/EXCCoin/exccwallet) servers
so users can purchase [proof-of-stake tickets](https://docs.excc.org/mining/proof-of-stake/)
on the [EXCCoin](https://excc.org/) network and have the pool of wallet servers
vote on their behalf when the ticket is selected.

## Architecture

![Stake Pool Architecture](https://i.imgur.com/2JDA9dl.png)

- It is highly recommended to use 3 exccd+exccwallet+stakepoold nodes for
  production use on mainnet.
- The architecture is subject to change in the future to lessen the dependence
  on exccwallet and MySQL.

## Git Tip Release notes

- The handling of tickets considered invalid because they pay too-low-of-a-fee
is now integrated directly into exccstakepool and stakepoold.
  - Users who pass both the adminIPs and the new adminUserIDs checks will see a new link on the
menu to the new administrative add tickets page.
  - Tickets are added to the MySQL database and then stakepoold is triggered to pull an update from the
database and reload its config.
  - To accommodate changes to the gRPC API, exccstakepool/stakepoold had
  their API versions changed to require/advertize 4.0.0. This requires
  performing the upgrade steps outlined below.
- **KNOWN ISSUE** Total tickets count reported by stakepoold may
  not be totally accurate until low fee tickets that have been added to
  the database can be marked as voted.  This will be resolved by future work. ([#201](https://github.com/EXCCoin/exccstakepool/issues/201)).

## Git Tip Upgrade Guide

1) Announce maintenance and shut down exccstakepool.
2) Upgrade Go to the latest stable version if necessary/possible.
3) Perform an upgrade of each stakepoold instance one at a time.
   * Stop stakepoold.
   * Build and restart stakepoold.
4) Edit exccstakepool.conf and set adminIPs/adminUserIDs appropriately to include
   the administrative staff for whom you wish give the ability to add low fee
   tickets for voting.
5) Upgrade and start exccstakepool after setting adminUserIDs.
6) Announce maintenance complete after verifying functionality.

## Requirements

- [Go](http://golang.org) 1.8.3 or newer.
- MySQL
- Nginx or other web server to proxy to exccstakepool

## Installation

#### Linux/BSD/MacOSX/POSIX - Build from Source

Building or updating from source requires the following build dependencies:

- **Go 1.8.3 or newer**

  Installation instructions can be found here: http://golang.org/doc/install.
  It is recommended to add `$GOPATH/bin` to your `PATH` at this point.

- **Dep**

  Dep is used to manage project dependencies and provide reproducible builds.
  To install:

  `go get -u github.com/golang/dep/cmd/dep`

Unfortunately, the use of `dep` prevents a handy tool such as `go get` from
automatically downloading, building, and installing the source in a single
command.  Instead, the latest project and dependency sources must be first
obtained manually with `git` and `dep`, and then `go` is used to build and
install the project.

- Run the following command to obtain the exccstakepool code and all dependencies:

```bash
$ git clone https://github.com/EXCCoin/exccstakepool $GOPATH/src/github.com/EXCCoin/exccstakepool
$ cd $GOPATH/src/github.com/EXCCoin/exccstakepool
$ dep ensure
```

- Assuming you have done the below configuration, build and run exccstakepool:

```bash
$ cd $GOPATH/src/github.com/EXCCoin/exccstakepool
$ go build
$ ./exccstakepool
```

- Build stakepoold and copy it to your voting nodes:

```bash
$ cd $GOPATH/src/github.com/EXCCoin/exccstakepool/backend/stakepoold
$ go build
```

## Updating

To update an existing source tree, pull the latest changes and install the
matching dependencies:

```bash
$ cd $GOPATH/src/github.com/EXCCoin/exccstakepool
$ git pull
$ dep ensure
$ go build
$ cd $GOPATH/src/github.com/EXCCoin/exccstakepool/backend/stakepoold
$ go build
```

## Setup

#### Pre-requisites

These instructions assume you are familiar with exccd/exccwallet.

- Create basic exccd/exccwallet/exccctl config files with usernames, passwords, rpclisten, and network set appropriately within them or run example commands with additional flags as necessary

- Build/install exccd and exccwallet from latest master

- Run exccd instances and let them fully sync

#### Stake pool fees/cold wallet

- Setup a new wallet for receiving payment for stake pool fees.  **This should be completely separate from the stake pool infrastructure.**

```bash
$ exccwallet --create
$ exccwallet
```

- Get the master pubkey for the account you wish to use. This will be needed to configure exccwallet and exccstakepool.

```bash
$ exccctl --wallet createnewaccount teststakepoolfees
$ exccctl --wallet getmasterpubkey teststakepoolfees
```

- Mark 10000 addresses in use for the account so the wallet will recognize transactions to those addresses. Fees from UserId 1 will go to address 1, UserId 2 to address 2, and so on.

```bash
$ exccctl --wallet accountsyncaddressindex teststakepoolfees 0 10000
```

#### Stake pool voting wallets

- Create the wallets.  All wallets should have the same seed.  **Backup the seed for disaster recovery!**

```bash
$ exccwallet --create
```

- Start a properly configured exccwallet and unlock it. See sample-exccwallet.conf.

```bash
$ exccwallet
```

- Get the master pubkey from the default account.  This will be used for votingwalletextpub in exccstakepool.conf.

```bash
$ exccctl --wallet getmasterpubkey default
```

#### MySQL

- Install, configure, and start MySQL

- Add stakepool user and create the stakepool database

```bash
$ mysql -uroot -ppassword

MySQL> CREATE USER 'stakepool'@'localhost' IDENTIFIED BY 'password';
MySQL> GRANT ALL PRIVILEGES ON *.* TO 'stakepool'@'localhost' WITH GRANT OPTION;
MySQL> FLUSH PRIVILEGES;
MySQL> CREATE DATABASE stakepool;
```

#### Nginx/web server

- Adapt sample-nginx.conf or setup a different web server in a proxy configuration

#### stakepoold
- Adapt sample-stakepoold.conf and run stakepoold.

#### exccstakepool

- Create the .exccstakepool directory and copy exccwallet certs to it
```bash
$ mkdir ~/.exccstakepool
$ cd ~/.exccstakepool
$ scp walletserver1:~/.exccwallet/rpc.cert wallet1.cert
$ scp walletserver2:~/.exccwallet/rpc.cert wallet2.cert
$ scp walletserver1:~/.stakepoold/rpc.cert stakepoold1.cert
$ scp walletserver2:~/.stakepoold/rpc.cert stakepoold2.cert
```

- Copy sample config and edit appropriately
```bash
$ cp -p sample-exccstakepool.conf exccstakepool.conf
```

## Running

The easiest way to run the stakepool code is to run it directly from the root of
the source tree:

```bash
$ cd $GOPATH/src/github.com/EXCCoin/exccstakepool
$ go build
$ ./exccstakepool
```

If you wish to run exccstakepool from a different directory you will need to change **publicpath** and **templatepath**
from their relative paths to an absolute path.

## Development

If you are modifying templates, sending the USR1 signal to the exccstakepool process will trigger a template reload.

## Operations

- exccstakepool will connect to the database or error out if it cannot do so

- exccstakepool will create the stakepool.Users table automatically if it doesn't exist

- exccstakepool attempts to connect to all of the wallet servers on startup or error out if it cannot do so

- exccstakepool takes a user's pubkey, validates it, calls getnewaddress on all the wallet servers, then createmultisig, and finally importscript.  If any of these RPCs fail or returns inconsistent results, the RPC client built-in to exccstakepool will shut down and will not operate until it has been restarted.  Wallets should be verified to be in sync before restarting.

- User API Tokens have an issuer field set to baseURL from the configuration file.
  Changing the baseURL requires all API Tokens to be re-generated.

## Adding Invalid Tickets

#### For Newer versions / git tip

If a user pays an incorrect fee, login as an account that meets the
adminUserIps and adminUserIds restrictions and click the 'Add Low Fee Tickets'
link in the menu.  You will be presented with a list of tickets that are
suitable for adding.  Check the appropriate one(s) and click the submit button.
Upon success, you should see the stakepoold logs reflect that the new tickets
were processed.

#### For v1.1.1 and below

If a user pays an incorrect fee you may add their tickets like so (requires exccd running with txindex=1):

```bash
exccctl --wallet stakepooluserinfo "MultiSigAddress" | grep -Pzo '(?<="invalid": \[)[^\]]*' | tr -d , | xargs -Itickethash exccctl --wallet getrawtransaction tickethash | xargs -Itickethex exccctl --wallet addticket "tickethex"
```

## Backups, monitoring, security considerations

- MySQL should be backed up often and regularly (probably at least hourly). Backups should be transferred off-site.  If using binary backups, do a test restore. For .sql files, verify visually.

- A monitoring system with alerting should be pointed at exccstakepool and tested/verified to be operating properly.  There is a hidden /status page which throws 500 if the RPC client is shutdown.  If your monitoring system supports it, add additional points of verification such as: checking that the /stats page loads and has expected information in it, create a test account and setup automated login testing, etc.

- Wallets should never be used for anything else (they should always have a balance of 0)

## Disaster Recovery

**Always keep at least one wallet voting while performing maintenance / restoration!**

- In the case of a total failure of a wallet server:
  * Restore the failed wallet(s) from seed
  * Restart the exccstakepool process to allow automatic syncing to occur.

## Issue Tracker

The [integrated github issue tracker](https://github.com/EXCCoin/exccstakepool/issues)
is used for this project.

## License

exccstakepool is licensed under the [copyfree](http://copyfree.org) MIT/X11 and
ISC Licenses.