"""Deployment scripts."""
import os

from fabric.api import cd, env, run, sudo
from fabric.colors import green, red
from fabric.contrib.files import exists
from fabric.contrib.project import rsync_project

env.use_ssh_config = True
install_dir = "/apps/goapi/"
project_path = "/go/src/github.com/etowett/go-api-sim/"
local_dir = "{0}{1}".format(os.getenv("HOME"), project_path)
live_dir = "/home/{0}/go/src/github.com/etowett/go-api-sim/"


def environment(arg):
    """Get the project environment."""
    if arg == "local":
        user = "vagrant"
    else:
        user = "focus"
    print((green("ENVIRONMENT USER: {0}".format(user))))
    return user


def dev():
    """Musale local."""
    # TODO: make this the default local host name
    env.hosts = ["stag"]


def stage():
    """Local host."""
    env.hosts = ["sms"]


def live():
    """Production host."""
    env.hosts = ["sdp"]


def deploy(arg=None):
    """Deployment script."""
    user = environment(arg)
    with cd(live_dir.format(user)):
        run("git pull origin master")
        run("go get")
        run("go build")
        run("go install")
    stop_goapi()
    with cd(install_dir):
        run("rm goapi")
        run("cp /home/{0}/go/bin/go-api-sim goapi".format(user))
    restart_goapi()
    return


def xdeploy(arg=None):
    """Local deployment script."""
    user = environment(arg)
    rsync_project(
        live_dir.format(user), local_dir=local_dir,
        exclude=["*.pyc", ".git*"], delete=True
    )
    with cd(live_dir.format(user)):
        print((green("build")))
        run("go build")
        print((green("install new")))
        run("go install")
    print((red("stop goapi")))
    stop_goapi()
    with cd(install_dir):
        print((red("remove old goapi and copy new")))
        run("rm goapi")
        run("cp /home/{0}/go/bin/go-api-sim goapi".format(user))
    print((green("start service")))
    restart_goapi()
    return


def setup(arg=None):
    """Get up."""
    user = environment(arg)
    sudo("yum -y install go git")
    if not exists("/home/{0}/go".format(user)):
        run("mkdir /home/{0}/go".format(user))
        run(
            "echo \"export GOPATH=$HOME/go\" "
            ">> /home/{0}/.bashrc".format(user))
    run("go get github.com/etowett/go-api-sim")
    with cd(live_dir.format(user)):
        run("go get")
        run("go build")
        run("go install")
    if not exists("/apps"):
        sudo("mkdir /apps")
        sudo("chown {0}:{1} /apps".format(user, user,))
    with cd("/apps"):
        if not exists("goapi"):
            run("mkdir goapi")
        with cd("goapi"):
            run("cp {0}env.sample .env".format(live_dir.format(user),))
            run("cp /home/{0}/go/bin/go-api-sim goapi".format(user))
    with cd("/var/log/"):
        if not exists("goapi"):
            sudo("mkdir goapi")
            sudo("chown {0}:{1} goapi".format(user, user,))
        with cd("goapi"):
            run("touch goapi.log")
    sudo(
        "cp {0}config/goapi.service "
        "/etc/systemd/system/goapi.service".format(live_dir.format(user),)
    )
    restart_goapi()
    return


def stop_goapi():
    """Stop the goapi service."""
    sudo("systemctl stop goapi")
    return


def restart_goapi():
    """Restart the goapi service."""
    sudo("systemctl restart goapi")
    return
