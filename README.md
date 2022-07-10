# argh command line parser

## background

The Go standard library [flag](https://pkg.go.dev/flag) way of doing things has long been
a source of frustration while implementing and maintaining the
[urfave/cli](https://github.com/urfave/cli) library. [Many alternate parsers
exist](https://github.com/avelino/awesome-go#standard-cli), including:

- [pflag](https://github.com/spf13/pflag)
- [argparse](https://github.com/akamensky/argparse)

In addition to these other implementations, I also got some help via [this
oldie](https://blog.gopheracademy.com/advent-2014/parsers-lexers/) and the Go standard
library [parser](https://pkg.go.dev/go/parser).

## goals

- get a better understanding of the whole problem space
- support both POSIX-y and Windows-y styles
- build a printable/JSON-able parse tree
- support rich error reporting

<!--
vim:tw=90
-->
