# bitwise-go

![Go](https://github.com/Akatsuki-py/bitwise-go/workflows/Go/badge.svg)

## Multi-platform terminal based bitwise calculator

_bitwise-go_ is multi base interactive calculator supporting dynamic base conversion and bit manipulation.
It's a handy tool for low level hackers, kernel developers and device drivers developers.

This repository is based on [_mellowcandle/bitwise_](https://github.com/mellowcandle/bitwise).

Some of the features include:
* Works on multi-platform such as Windows, OSX, Linux.
* Command line calculator supporting all bitwise operations.
* Individual bit manipulator.

![conversion](https://imgur.com/pcth8U0.png "Bitwise conversion")

![interactive](https://imgur.com/QI9BrHl.png "Bitwise interactive")

## Install

#### Download Binary

Please download binary from [Release page](https://github.com/Akatsuki-py/bitwise-go/releases).

#### Build by yourself

Requirement: Go and Make

```sh
$ git clone https://github.com/Akatsuki-py/bitwise-go.git
$ cd ./bitwise-go
$ make
```

## Usage
_bitwise-go_ can be used both Interactively and in command line mode.

### Command line calculator mode
In command line mode, bitwise will calculate the given expression and will output the result in all bases including binary representation.

_bitwise-go_ detects the base by the preface of the input (_0x/0X_ for hexadecimal, leading _0_ for octal, _0b_ for binary, and the rest is decimal).

### Interactive mode
_bitwise-go_ starts in interactive mode if no command line parameters are passed. In this mode, you can input a number and manipulate it and see the other bases change dynamically.
It also allows changing individual bits in the binary.

#### Navigation in interactive mode
To move around use the arrow keys, or use _vi_ key bindings : <kbd> h </kbd> <kbd> j </kbd> <kbd> k </kbd> <kbd> l </kbd>.
Leave the program by pressing <kbd> q </kbd>.

##### Binary specific movement
You can toggle a bit bit using the <kbd> space </kbd> key.

#### others
* _q_, _esc_ - Exit
