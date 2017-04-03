#!/usr/bin/env python

from fabric.api import env, cd, run, sudo

env.use_ssh_config = True
env.hosts = ["mpesa"]
code_dir = "/home/ekt/go/src/github.com/etowett/go-api-sim/"
install_dir = "/apps/goapi/"


def deploy():
    with cd(code_dir):
        run("git pull origin master")
        run("go build")
        run("go install")
    sudo("systemctl stop goapi")
    with cd(install_dir):
        run("rm goapi")
        run("cp /home/ekt/go/bin/go-api-sim .")
    sudo("systemctl start goapi")
    return
