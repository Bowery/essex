from fabric.api import *
import requests

project = "mercer"
repository = "git@github.com:Bowery/" + project + ".git"
hosts = [
  'ubuntu@ec2-54-146-175-128.compute-1.amazonaws.com'
]
env.key_filename = '/home/ubuntu/.ssh/id_aws'
env.password = 'java$cript'

@parallel
def restart():
  sudo('mkdir -p /home/ubuntu/gocode/src/github.com/Bowery/')
  with cd('/home/ubuntu/gocode/src/github.com/Bowery/gopackages'):
    run('git pull')

  with cd('/home/ubuntu/gocode/src/github.com/Bowery/' + project):
    run('git pull')
    
    sudo('GOPATH=/home/ubuntu/gocode go get -d')
    sudo('GOPATH=/home/ubuntu/gocode go build')

    sudo('cp -f ' + project + '.conf /etc/init/' + project + '.conf')
    sudo('initctl reload-configuration')
    sudo('restart ' + project)

def deploy():
  execute(restart, hosts=hosts)
