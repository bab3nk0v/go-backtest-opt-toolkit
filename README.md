# Intro

This research project has been developed during my master's thesis preparation.
Includes all the code and data needed to run your own experiments.

# Quickstart

## Database
first of all we need mysql db to store intermediate results\
the simplest way is to set it up via docker:

```console
$ docker pull mysql:8.0
```

```console
$ docker run \
-d \
--rm \
-p 3306:3306 \
-e MYSQL_USER=goptuna \
-e MYSQL_DATABASE=goptuna \
-e MYSQL_PASSWORD=password \
-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
--name goptuna-mysql \
mysql:8.0
```

## Building

### Install ta-lib

#### Linux

```console
$ wget http://prdownloads.sourceforge.net/ta-lib/ta-lib-0.4.0-src.tar.gz
$ tar -xzf ta-lib-0.4.0-src.tar.gz
$ cd ta-lib
$ ./configure --prefix=/usr LDFLAGS="-lm"
$ make
$ sudo make install
```

#### MacOS

```console
$ brew install ta-lib
```

### Compile optimizer

To compile we need golang compiler of version at least 1.15

```console
$ git clone https://github.com/ababenkov/trade-optimizer
$ cd trade-optimizer
$ go build .
```

Go build outputs trade-optimizer binary for your platform

### Run!

Running the binary will run the optimizer with default settings
```console
$ ./trade-optimizer
```

Custom options can be also specified

```console
 $ ./trade-optimizer --help
Usage of ./trade-optimizer:
  -dataPath string
        path to OHLCV dataset (default "datasets/XMR_1m.csv")
  -dbHost string
        db host (default "localhost")
  -dbLogin string
        db username (default "goptuna")
  -dbName string
        database name (default "goptuna")
  -dbPass string
        db pass (default "password")
  -dbPort int
        db port (default 3306)
  -fee float
        fee
  -nItersNoChange int
        number of iterations without improvement to stop (default 3000)
  -pairname string
        pair name (default "XMR_USD")
  -resultFolder string
        where to dump optimization results (default "opt_results")
  -timeInMillseconds
        true=time in milliseconds false=time in seconds (default true)
  -tradeSize float
        trade size (default 100)

```
The result of running the optimizer is a json file with 
parameters and trades. The analysis is done in analysis.ipynb

