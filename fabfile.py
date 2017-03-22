"""Fabfile to deploy stuff."""

from fabric.api import env, cd, run, sudo, local, lcd, put
from fabric.contrib.files import exists

import os


app_dir = '/apps/smsl'
git_repo = 'git@bitbucket.org:teamictlife/smsleopard.git'

tmp = "/tmp/sms"
tmp_f = "%s/smsl.tar.gz" % tmp

env.use_ssh_config = True
env.hosts = ['smsl']


def xdeploy():
    if os.path.exists(tmp):
        local('rm -rf %s' % tmp)
    local('mkdir %s' % tmp)
    with lcd(app_dir):
        local('tar -czhf %s smsleopard --exclude=".git*"' % (tmp_f))
    if exists(tmp):
        run('rm -rf %s' % tmp)
    run('mkdir %s' % tmp)
    put(tmp_f, tmp_f)
    with cd(app_dir):
        if exists('smsleopard'):
            run('rm -rf smsleopard')
        run('tar -xzf %s' % tmp_f)
        with cd('/var/log'):
            if not exists('smsleopard'):
                sudo('mkdir -p smsleopard/app')
                sudo('chown -R %s:%s smsleopard' % (user, user,))
                run('touch smsleopard/app/smsleopard.log')
    prep_remote()
    restart_services()
