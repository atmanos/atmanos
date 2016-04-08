# Developing AtmanOS

## Build AtmanOS

If you've just downloaded the AtmanOS project, run `bin/setup` to download the
required dependencies and build AtmanOS.

Once you've run `bin/setup`, rebuilding AtmanOS can be done with `make build`.

## Testing changes

You can test changes using basic workflow outlined in [Running locally with
Vagrant](./running-locally-with-vagrant.md).

There's also a helper script which automates the process of rebuilding AtmanOS
and deploying a test program:

```
$ bin/rebuild-and-deploy github.com/atmanos/example/hello
```

It may be useful to automatically build and deploy whenever changes are made,
which can be done with [entr](http://entrproject.org/):

```
$ find src -type f |
  entr bin/rebuild-and-deploy github.com/atmanos/example/hello
```

## Debug an AtmanOS kernel

You can get a GDB session for a running (or crashed) kernel with:

```
$ vagrant ssh -- debugvm hello
```
