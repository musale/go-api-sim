
import os

from fabric.api import env, cd, run, sudo, local, lcd, put
from fabric.contrib.files import exists

env.use_ssh_config = True
env.hosts = ["sms"]

code_dir = "/home/ekt/go/src/github.com/etowett/go-api-sim/"
install_dir = "/apps/goapi/"
local_dir = "/home/ekt/go/src/github.com/etowett/"
live_dir = "/home/focus/go/src/github.com/etowett/"
tmp = "/tmp/goapi"
tmp_f = "%s/goapi.tar.gz" % tmp
user = "focus"


def deploy():
    with cd("%s/go-api-sim/" % (live_dir,)):
        run("git pull origin master")
        run("go get")
        run("go build")
        run("go install")
    stop_goapi()
    with cd(install_dir):
        run("rm goapi")
        run("cp /home/ekt/go/bin/go-api-sim goapi")
    restart_goapi()
    return


def xdeploy():
    if os.path.exists(tmp):
        local('rm -rf %s' % tmp)
    local('mkdir %s' % tmp)
    with lcd(local_dir):
        local('tar -czhf %s go-api-sim --exclude=".git*"' % (tmp_f))
    if exists(tmp):
        run('rm -rf %s' % tmp)
    run('mkdir %s' % tmp)
    put(tmp_f, tmp_f)
    with cd(live_dir):
        if exists('go-api-sim'):
            run('rm -rf go-api-sim')
        run('tar -xzf %s' % tmp_f)
    with cd('%s/go-api-sim' % live_dir):
        run('go get')
        run('go build')
        run('go install')
    restart_goapi()
    return


def setup():
    sudo("yum -y install go git")
    if not exists("/home/focus/go"):
        run("mkdir /home/focus/go")
        run("echo \"export GOPATH=$HOME/go\" >> /home/focus/.bashrc")
    run("go get github.com/etowett/go-api-sim")
    with cd('%sgo-api-sim' % live_dir):
        run('go get')
        run('go build')
        run('go install')
    if not exists("/apps"):
        sudo("mkdir /apps")
        sudo("chown %s:%s /apps" % (user, user,))
    with cd("/apps"):
        if not exists("goapi"):
            run("mkdir goapi")
        with cd("goapi"):
            run("cp %sgo-api-sim/env.sample .env" % (live_dir,))
            run("cp /home/focus/go/bin/go-api-sim goapi")
    with cd("/var/log/"):
        if not exists("goapi"):
            sudo("mkdir goapi")
            sudo("chown %s:%s goapi" % (user, user,))
        with cd("goapi"):
            run("touch goapi.log")
    sudo(
        "cp %sgo-api-sim/config/goapi.service "
        "/etc/systemd/system/goapi.service" % (live_dir,)
    )
    restart_goapi()
    return


def stop_goapi():
    sudo("systemctl stop goapi")
    return


def restart_goapi():
    sudo('systemctl restart goapi')
    return
