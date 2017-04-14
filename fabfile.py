
from fabric.colors import green, red
from fabric.contrib.files import exists
from fabric.api import env, cd, run, sudo
from fabric.contrib.project import rsync_project

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
    rsync_project(
        live_dir, local_dir='%sgo-api-sim' % local_dir,
        exclude=['*.pyc', '.git*'], delete=True
    )
    with cd('%sgo-api-sim' % live_dir):
        print(green("get dependencies if any"))
        run('go get')
        print(green("build"))
        run('go build')
        print(green("install new"))
        run('go install')
    print(red("stop goapi"))
    stop_goapi()
    with cd(install_dir):
        if exists("goapi"):
            print(red("remove old goapi"))
            run("rm goapi")
        print(green("copy new goapi"))
        run("cp /home/focus/go/bin/go-api-sim goapi")
    print(green("start service"))
    restart_goapi()
    return


def setup():
    sudo("yum -y install go git")
    if not exists("/home/focus/go"):
        run("mkdir /home/focus/go")
        run("echo \"export GOPATH=$HOME/go\" >> /home/focus/.bashrc")
    rsync_project(
        live_dir, local_dir='%sgo-api-sim' % local_dir,
        exclude=['*.pyc', '.git*'], delete=True
    )
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
