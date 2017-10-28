> *NOTE*: AtmanOS is an unstable and experimental project!

# AtmanOS

AtmanOS allows you to compile ordinary Go programs into standalone unikernels
that run under the Xen hypervisor.

AtmanOS is implemented as a series of patches and additional files for Go's
runtime and standard library. AtmanOS is a `GOOS` available from the modified
`go` command, meaning programs can be cross-compiled for AtmanOS in the normal
manner:

```
GOOS=atman go build
```

## Build AtmanOS

Build AtmanOS by running `bin/setup`,
which will download the required dependencies to the build directory
and then build AtmanOS itself.

## Build a program and deploy it

Read the [Running locally with Vagrant](doc/running-locally-with-vagrant.md)
tutorial for the fastest way to build your first program with AtmanOS.

## Supported I/O

  * _Console_ - available on os.Stdin and os.Stdout
  * _Network_ - managed out-of-tree in [atmanos/anet](https://github.com/atmanos/anet)

Contributing
------------

We love pull requests from everyone.
By participating in this project,
you agree to abide by the [Go Community Code of Conduct][code of conduct].

[code of conduct]: https://golang.org/conduct

We expect everyone to follow the code of conduct
anywhere in the project's codebases,
issue trackers, chatrooms, and mailing lists.

License
-------

AtmanOS is Copyright (c) 2015 Bernerd Schaefer. It is free software,
and may be redistributed under the terms specified in the [LICENSE] file.

[LICENSE]: /LICENSE
