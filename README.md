# Do..Emm...

[![Build Status](https://travis-ci.org/raomuyang/doemm.svg?branch=master)](https://travis-ci.org/raomuyang/doemm)
[![Go Report Card](https://goreportcard.com/badge/github.com/raomuyang/doemm)](https://goreportcard.com/report/github.com/raomuyang/doemm)
[![GoDoc](https://godoc.org/github.com/raomuyang/doemm?status.svg)](https://godoc.org/github.com/raomuyang/doemm)

```
______         _____
|  _  \       |  ___|
| | | |___    | |__ _ __ ___  _ __ ___
| | | / _ \   |  __| '_ ` _ \| '_ ` _ \
| |/ / (_) | _| |__| | | | | | | | | | |
|___/ \___(_|_)____/_| |_| |_|_| |_| |_|
```


A handy command line tool to help you manager complicated commands.

## Quick start

Install via golang or download from [releases](https://github.com/raomuyang/doemm/releases)
```shell
go get github.com/raomuyang/doemm
```

Record complex operations with simple commands
> supports encrypted text save mode if you need


1. config for sync with gist, encrypt the information text of command lines to saved in application home (~/.doemm)

``` shell
doemm config -gist <github-oAuth-token> [-encrypt]
```

2. save a command and give it an alias
```shell
# doeem -alias <command alias> <target command to save>
doemm -alias db-host ssh -p 22 user@fake.machine.balabala
```

3. quick call
```shell
doemm -s db-host
```

4. multi-line input

```shell
doemm
```

then you will enter the interactive interface
```
Input multipart commands in interactive
1. type emm.done to exit
2. type emm.alias.ALIAS_NAME to set alias of commands
emm...>
```

5. push/pull to gist
> If you save sensitive information, it is highly recommended that you do not synchronize or synchronize separately with the -single parameter (although you are creating a secret gist)

```
doemm pull/push [-single <command-alias>]
```

6. other operations

* list all

```shell
doemm -list
```

* print commond(s) info with specified alias

```shell
doemm -print <command alias>
```

* delete
> -sync: synchronous delete operation into gist

```shell
doemm rm -t <target-alias> [-sync]
```

