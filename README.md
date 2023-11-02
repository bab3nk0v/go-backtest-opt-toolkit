# 🚀 Welcome to the High-Performance Trading Strategy Optimization Toolkit

## Introduction 📘

Welcome, fellow data enthusiast! You've stumbled upon a repository that served as a cornerstone for my academic pursuit during my master's thesis. This project isn't just a collection of code and data; it's a testament to the fusion of cutting-edge optimization techniques and high-performance computing in Go. Whether you're an aspiring quantitative researcher, a data science maven, or simply a curious mind, you'll find this toolkit to be a robust playground for running extensive trading strategy experiments.

**Disclaimer:** While this project has been crafted with meticulous attention to detail, it is currently not actively maintained and was developed primarily for educational purposes.

## Getting Started 🚦

### Prerequisites: Setting up the Database 🛢️

Our journey begins with setting up a MySQL database to store the fruits of our optimization labor. But fear not, for Docker comes to the rescue with a swift setup:

```console
$ docker pull mysql:8.0
```

Proceed by conjuring up your SQL instance with this incantation:

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

### Compilation Cruise: Install & Build 🛠️

#### Technical Alchemy on Linux

```console
$ wget http://prdownloads.sourceforge.net/ta-lib/ta-lib-0.4.0-src.tar.gz
$ tar -xzf ta-lib-0.4.0-src.tar.gz
$ cd ta-lib
$ ./configure --prefix=/usr LDFLAGS="-lm"
$ make
$ sudo make install
```

#### For the Mac OS Magicians

```console
$ brew install ta-lib
```

#### Conjure the Optimizer

Your quest for compiling this project begins with the mighty Golang compiler (v1.15+):

```console
$ git clone https://github.com/ababenkov/trade-optimizer
$ cd trade-optimizer
$ go build .
```

Voilà! The `trade-optimizer` binary tailored for your platform is ready to unveil its magic.

### Launch Instructions 🚀

Let the binary dance to the default tune:

```console
$ ./trade-optimizer
```

Or, conduct the performance with your own symphony of options:

```console
 $ ./trade-optimizer --help
```

### The Grand Finale: Results 🎉

Executing the optimizer culminates in a JSON opus, detailing the parameters and trades of your strategic masterpiece. All of which are ready for your analytical concerto in `analysis.ipynb`.

---

Embark on this journey to explore the symphony of algorithms and performance, and may your data be ever in your favor. Enjoy, and may the odds bring forth strategies of great prosperity!
