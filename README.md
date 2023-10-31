# v
A version manager you might not want to use.

## Overview

`v` is a simple version manager inspired from other tools like [asdf](https://github.com/asdf-vm/asdf), [pyenv](https://github.com/pyenv/pyenv), [n](https://github.com/tj/n) and [nvm](https://github.com/nvm-sh/nvm). At it's core, it's a reinvention of the wheel with some extras.

- First and foremost, while the first version is about Python version management, the plan is to expand to support a bunch more runtime (with an emphasis on simplifying adding more runtimes to manage);
- A lot of those tools are written as shellscript, which I find somewhat inscrutable. Go is a bit easier to read;
- ...? It's a reason to write some Go. :)
