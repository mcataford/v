# v
A version manager you might not want to use.

> # ✈️ Moved away!
>
> This project has moved away from Github and is now hosted [elsewhere](https://forge.karnov.club/marc/v).

## Overview

`v` is a simple version manager inspired from other tools like [asdf](https://github.com/asdf-vm/asdf), [pyenv](https://github.com/pyenv/pyenv), [n](https://github.com/tj/n) and [nvm](https://github.com/nvm-sh/nvm). At it's core, it's a reinvention of the wheel with some extras.

- First and foremost, while the first version is about Python version management, the plan is to expand to support a bunch more runtime (with an emphasis on simplifying adding more runtimes to manage);
- A lot of those tools are written as shellscript, which I find somewhat inscrutable. Go is a bit easier to read;
- ...? It's a reason to write some Go. :)

## Roadmap

While the plan for the first release is to only support Python runtimes, expanding to others will be next so that `v` can just handle all/most version management needs.

## Usage

### Building your own and setting up

Pre-built binaries are not currently available. You can clone the repository and build your own via `. scripts/build`.

You should find a suitable place for the binary (`/usr/local/bin` is a good location) and if not already included, add its location to `$PATH`.

Finally, run `v init` to create directories to store artifacts and state (under `~/.v` unless override using the
`V_ROOT` environment variable). The following should also be added to your shell's configuration (i.e. `.zshrc`,
`.bashrc`, ...):

```sh
export PATH=<path-to-v-executable>:$PATH
eval "$(v init --add-path)"
```

This will handle adding shim paths to your shell without hassle.

### Usage

`v` will print a helpful list of available commands.

The most important things to know include `v python install <version>` to install new versions and `v python use <installed version>` to use a specific version of Python.

## Contributing

The project isn't currently accepting contributions because it's not yet set up to do so. Stay tuned.
