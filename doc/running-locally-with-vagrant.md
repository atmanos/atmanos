# Running locally with Vagrant

The AtmanOS includes a Vagrantfile whch describes a development Xen
environment. To use this development environment, you'll need:

  * Vagrant: https://www.vagrantup.com/
  * VirtualBox: https://www.virtualbox.org/

If you're on OS X and use Homebrew, you can get these with:

```
$ brew cask install vagrant virtualbox
```

## Build AtmanOS

If you've just downloaded the AtmanOS project, run `bin/setup` to download the
required dependencies and build AtmanOS.

Once you've run `bin/setup`, rebuilding AtmanOS can be done with `make build`.

## Provision the virtual machine

From the root directory of the AtmanOS project, run:

```
$ vagrant up
```

The first time you run this, it will take a while, as it needs to install Xen
and reboot the virtual machine.

Once the command completes, you should be able to connect to the machine by
running `vagrant ssh`.

At any time you can use `vagrant halt` to shut down the virtual machine, or
`vagrant pause` to pause it for faster boot later.

## Run an AtmanOS kernel

For this section, we'll be using the hello program from the [repository of
AtmanOS example programs][example].

  [example]: github.com/atmanos/example

Download the program with:

```
$ go get github.com/atmanos/example/hello
```

You can run `hello` locally to see that it first prints "Hello, world", and
then prints the current time every few seconds.

But we want to run the `hello` program on Xen!

To do that, we'll first need to build a kernel image with `atman`:

```
$ bin/atman build -o vagrant/images/hello \
  github.com/atmanos/example/hello
```

We told `atman` to put the built image in `vagrant/images`. That allows up to
transfer the image to the Vagrant environment with:

```
$ vagrant rsync
```

Finally, we can run our `hello` kernel with:

```
$ vagrant ssh -- startvm hello
```

The [`startvm` command](../vagrant/bin/startvm) is a wrapper for a few Xen
commands, stopping any existing machine with the same name, and then starting a
new machine with the [default template](../vagrant/template.xl).

To see the output of `hello`, attach to the console with [another Xen
wrapper](../vagrant/bin/console):

```
$ vagrant ssh -- console hello
Hello, world
The current time is 2016-04-23 23:53:05.651799283 +0000 UTC
The current time is 2016-04-23 23:53:10.651806585 +0000 UTC
The current time is 2016-04-23 23:53:15.651807123 +0000 UTC
```

Success!

Now we can stop the `hello` machine with:

```
$ vagrant ssh -- sudo xl destroy hello
```
