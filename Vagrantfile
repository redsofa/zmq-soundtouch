# -*- mode: ruby -*-
# vi: set ft=ruby :
Vagrant.configure(2) do |config|

  config.vm.box = "geerlingguy/centos7"
  
  config.vm.provider "virtualbox" do |vb|
    vb.name = "zmqsoundtouch"
  end
 
  config.vm.provision "install ansible", type: "shell",
    inline: "yum -y install ansible"
  
  # Share a folder we can use for coding
  config.vm.synced_folder "src", "/usr/local/src", disabled: false
  
  # run ansible as shell script instead of dedicate provisioner,
  # because we're waiting for Vagrant 1.8.2 (which will make ansible_local work correctly)
  # check this bug for update: https://github.com/geerlingguy/drupal-vm/issues/450
  # or just check periodically for Vagrant 1.8.2 to be released:
  # https://www.vagrantup.com/downloads.html
  # and then switch to ansible_local provisioner:
  # https://www.vagrantup.com/docs/provisioning/ansible_local.html
  # This will show all the ansible status in command line
  config.vm.provision "ansible-playbook", type: "shell",
    inline: "ansible-playbook /vagrant/playbook.yml"
 
end