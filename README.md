<!-- Copyright 2022, 2025 The Autogun Authors. All rights reserved.
     Use of this source code is governed by a BSD-style
     license that can be found in the LICENSE file. 
-->

<h1 align="center">
  <br>
    AUTOGUN
  <br>
  <br>
</h1>

<h4 align="center">Headless browser automation with reusable workflows </h4>

**Autogun** is a configuration language and toolchain to automate headless browsers.  Currently, it uses the [Chrome DevTools protocol][] to automate [Chromium][]-based browsers.

> [!important] Autogun is pre-alpha software. Features are subject to change, and compatibility may vary until Autogun reaches 1.0.

## Installation

To get started with Autogun, install the executable using [Go][]:

```shell
go install github.com/Carbonfrost/autogun/cmd/autogun@latest
```

then use `autogun` from your scripts. For example, to navigate to a webpage and obtain its title:

```shell
autogun run https://github.com/Carbonfrost/autogun -title
```

Check out [our docs][] for more guidance on how to use Autogun.

### Development

1. Download and install [Go 1.24+](https://go.dev)
2. Run `go build ./...`
3. Run the tests using `go test ./...`

## License

This project uses a [BSD-style license](LICENSE).

[Chrome DevTools protocol]: https://chromedevtools.github.io/devtools-protocol/
[Chromium]: https://www.chromium.org/Home/
[Go]: https://go.dev
[our docs]: https://github.com/Carbonfrost/autogun
