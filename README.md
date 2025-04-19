# spin

[![License](https://img.shields.io/github/license/FollowTheProcess/spin)](https://github.com/FollowTheProcess/spin)
[![Go Reference](https://pkg.go.dev/badge/github.com/FollowTheProcess/spin.svg)](https://pkg.go.dev/github.com/FollowTheProcess/spin)
[![Go Report Card](https://goreportcard.com/badge/github.com/FollowTheProcess/spin)](https://goreportcard.com/report/github.com/FollowTheProcess/spin)
[![GitHub](https://img.shields.io/github/v/release/FollowTheProcess/spin?logo=github&sort=semver)](https://github.com/FollowTheProcess/spin)
[![CI](https://github.com/FollowTheProcess/spin/workflows/CI/badge.svg)](https://github.com/FollowTheProcess/spin/actions?query=workflow%3ACI)

A very simple terminal spinner

<p align="center">
<img src="https://github.com/FollowTheProcess/spin/raw/main/docs/img/demo.gif" alt="demo">
</p>

## Project Description

I needed a very simple, minimal overhead, nice looking terminal spinner that didn't bring in a shed load of dependencies. So here it is!

It has a few nicities:

- Auto terminal detection and colouring via [hue]
- Customisable colours
- Custom progress message
- Simple and convenient API

## Installation

```shell
go get github.com/FollowTheProcess/spin@latest
```

## Quickstart

```go
package main

import (
    "os"
    "time"
    
    "github.com/FollowTheProcess/spin"
)

func main() {
    spinner := spin.New(os.Stdout, "Digesting")

    spinner.Start()
    defer spinner.Stop()

    time.Sleep(2 * time.Second)
}
```

You can also wrap a function in a spinner...

```go
package main

import (
    "os"
    "time"
    
    "github.com/FollowTheProcess/spin"
)

func main() {
    spinner := spin.New(os.Stdout, "Digesting")

    // This is equivalent to the example above, Do will handle
    // starting and stopping the spinner for you
    spinner.Do(func() {
        time.Sleep(2 * time.Second)
    })
}
```

### Credits

This package was created with [copier] and the [FollowTheProcess/go_copier] project template.

[copier]: https://copier.readthedocs.io/en/stable/
[FollowTheProcess/go_copier]: https://github.com/FollowTheProcess/go_copier
[hue]: https://github.com/FollowTheProcess/hue
