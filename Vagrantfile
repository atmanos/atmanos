# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.configure(2) do |config|
  config.vm.box = "ubuntu/trusty64"

  # Disable the default shared folder, because we can't use VirtualBox shared
  # folders when running under Xen.
  config.vm.synced_folder ".", "/vagrant", disabled: true

  # Run vagrant rsync to update the shared images.
  config.vm.synced_folder "vagrant/", "/home/vagrant/atman", type: "rsync"

  # Install xen hypervisor
  config.vm.provision "shell", inline: <<-SHELL
    sudo apt-get update
    sudo apt-get install -y xen-hypervisor-4.4-amd64 gdb
    sed -i '1s|^|export PATH=$HOME/atman/bin:$PATH\\n|' .bashrc

    sudo ln -nsf /home/vagrant/atman/net/xenbr0.cfg /etc/network/interfaces.d/xenbr0.cfg
    sudo ln -nsf /home/vagrant/atman/net/xenbr0-up /etc/network/if-up.d/xenbr0-up

    sudo reboot
  SHELL
end
