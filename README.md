# Simple-Password-Generator

Creation of a simple password generator usable in command line.

## Table of contents

1. [Simple run](#simple-run)
2. [Create an executable program](#executable-binary-program)
3. [Usage](#usage)

**N.B.:** To use all the command, you must have installed the [Golang environment](https://golang.org/).

## Simple run

You can run the program directly with :

```shell
$ go run passwordgenerator.go
```

## Executable binary program

You can run the followig command :

```shell
$ go build passwordgenerator.go
```

## Usage

You can use the program from 2 ways :

- you can just open the binary file which open an interactive program;
- you can just pass the arguments to the command when you call it :

```shell
$ passwordgenerator.exe <length:int> <number_of_digits:int> <number_of_symbols:int> <allow_uppercase:(false|true)> <allow_repeat:(false|true)>
```
