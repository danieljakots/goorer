# Goorer

Goorer is a personal finance software. It's loosely based on
[poorer](https://framagit.org/Steap/poorer). Goorer is written in go — hence
the name, and has a CLI instead of a web one.

The reason for this rewrite is I got tired of maintaining my own poorer fork.
Upstream didn't care about my bug fixes and misc improvements. I had to use a
virtual env which frequently stop working for no reason.

# Using goorer

## Installation

Get the repository and run `make build`. It will produce the *goorer* binary.

## Goorer's data

Goorer gets its data from yaml files. There are two types of files:
- categories
- money exchanges

The *categories.yml* file is a key value file which is used only for spendings.
The key is who got your money, and the value is its categories.

An example could be:
~~~
rent: home
furniture shop: home
supermarket: food
restaurant: food
~~~

All entities need to have one category.

Money exchanges file are a file with the earnings and spendings listed. For
instance:
~~~
earnings:
- date: 2020-12-25
  amount: 4321
  with: Company
- date: 2020-12-25
  amount: 5
  with: Santa Claus

spendings:
- date: 2020-12-01
  amount: 1234
  with: rent
- date: 2020-12-12
  amount: 13.37
  with: cat food shop
- date: 2020-12-21
  amount: 73.31
  with: saq
- date: 2020-12-25
  amount: 42.24
  with: cat food shop
~~~

Poorer reads only files with the `.yaml` extension. To allow have both file in
the same directory, goorer reads only files with the `.yml` extension since the
format changes a bit. With sed(1) you can easily move from poorer to goorer,
you just need to search and replace `to:` and `from:` to `with:`.

## Usage

Here's how to use goorer:

~~~
$ goorer -h
usage: goorer [-h] {summary, earnings, spendings} path/to/data

Each subcommand accepts a --date YYYY[-MM] to filter on a subset of entries
The spendings subcommand accept a -d/--details. This prints spendings
without using categories.
~~~

Here's what it looks like:

(testdata is the directory with some data for unit tests).

### Summary

~~~
$ goorer summary testdata
You earned $5448.00
You spent $3643.88
You saved $1804.12
You spent 66.88% of your earnings
~~~

### Earnings

~~~
$ goorer earnings testdata
From Company                  : we earned $5443.00,   this is 99.91%
From Santa Claus              : we earned $5.00,      this is  0.09%
~~~

### Spendings

~~~
$ goorer spendings testdata
For home                     : we spent $3445.00,   this is 94.54%
For cat                      : we spent $125.57,    this is  3.45%
For wine                     : we spent $73.31,     this is  2.01%
~~~

### Filter by date

~~~
$ goorer earnings --date 2019 testdata
From Company                  : we earned $1122.00,   this is 100.00%
~~~

### Don't use categories to print spendings

~~~
$ goorer spendings --details testdata
For rent                     : we spent $3445.00,   this is 94.54%
For cat food shop            : we spent $125.57,    this is  3.45%
For saq                      : we spent $73.31,     this is  2.01%
~~~
