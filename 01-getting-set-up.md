# 01 - Getting set up

We'll start by getting the Hashicorp stack set up on your machine. We'll be using one machine for each team, so you don't have to set this up on all of them.

## Installing the Hashicorp stack

The key components that we'll be using are [Nomad](https://nomadproject.io) and [Consul](https://consul.io).

Nomad is an application scheduling tool which allows "jobs" to be submitted for execution on a cluster of machines, and handles things like running multiple instances and recovering without downtime if a machine fails.

Consul is a service discovery, configuration management and service mesh tool. Nomad uses Consul to allow applications to find other applications and services that are running in the cluster.

They're both distributed as single binaries, so they're pretty straightforward to install.

On a Mac, you can run the following script to download both applications and install them into `/usr/local/bin`:

```bash 
curl https://releases.hashicorp.com/nomad/0.8.4/nomad_0.8.4_darwin_amd64.zip -o /tmp/nomad.zip && 
curl https://releases.hashicorp.com/consul/1.2.2/consul_1.2.2_darwin_amd64.zip -o /tmp/consul.zip &&
unzip /tmp/nomad.zip -d /usr/local/bin && unzip /tmp/consul.zip -d /usr/local/bin
```

If you're on Linux, the script should be roughly the same with different binaries, but you might have to tweak it depending on your distribution - all the commands need to be run as root:

```bash
curl https://releases.hashicorp.com/nomad/0.8.4/nomad_0.8.4_linux_amd64.zip -o /tmp/nomad.zip && 
curl https://releases.hashicorp.com/consul/1.2.2/consul_1.2.2_linux_amd64.zip -o /tmp/consul.zip &&
unzip /tmp/nomad.zip -d /usr/local/bin && unzip /tmp/consul.zip -d /usr/local/bin
```

On Linux you'll want to also make sure that the `/usr/local/bin` path is in your shell's PATH variable and add it if it is not:
```bash
root@balloon# echo $PATH
/usr/local/sbin:/usr/sbin:/usr/bin:/sbin:/bin:/snap/bin
root@balloon# export PATH=$PATH:/usr/local/bin
```

Then you can check that everything installed okay:

```
→ nomad version && consul version
Nomad v0.8.4 (dbee1d7d051619e90a809c23cf7e55750900742a)
Consul v1.2.2
Protocol 2 spoken by default, understands 2 to 3 (agent will automatically use protocol >2 when speaking to compatible agents)
```

That's all we need to get the Hashicorp compoonents installed!


## Installing Docker

Nomad is capable of running a variety of different kinds of jobs. It can run an executable file that's installed on a machine, spin up virtual machines to run each job, or even use QEMU to run binaries for different architectures.

To keep things nice and clean and isolated, we will be using the Docker driver. This allows Nomad to run a job that's entirely contained inside a Docker image. If you're not familiar with Docker, you can think of it as a mechanism that that allows us to run "images" containing an application and all of its dependencies in an isolated, custom environment called a "container" – kind of like a really light-weight virtual machine.

To install Docker on a Mac, grab the image from [https://download.docker.com/mac/stable/Docker.dmg](https://download.docker.com/mac/stable/Docker.dmg).

To install it on Linux, you'll have to refer to your distribution's instructions ([Ubuntu guide](https://docs.docker.com/install/linux/docker-ce/ubuntu/#install-using-the-repository))

Check that it's installed okay:

```
→ docker version
Client:
 Version:      18.03.1-ce
 API version:  1.37
 Go version:   go1.9.5
 Git commit:   9ee9f40
 Built:        Thu Apr 26 07:13:02 2018
 OS/Arch:      darwin/amd64
 Experimental: false
 Orchestrator: swarm

Server:
 Engine:
  Version:      18.03.1-ce
  API version:  1.37 (minimum version 1.12)
  Go version:   go1.9.5
  Git commit:   9ee9f40
  Built:        Thu Apr 26 07:22:38 2018
  OS/Arch:      linux/amd64
  Experimental: true
```

You should see some output roughly simiar to the above, but you don't neesd to worry about the version too much!


## Running the Hashicorp stack

Okay, so now we've got the Hashicorp stack installed. Let's run a cluster!

We'll start a Consul agent first. We'll be running it in development mode, which automates some settings that we don't have to worry about for now.

You can start Consul in the foreground using the following command:

```bash
consul agent --data-dir=/tmp/consul -dev
```

Consul should start in the foreground and output logging messages with useful debugging information if anything goes wrong.

Next we can run Nomad in a new terminal:

```bash
nomad agent -client -server -data-dir=/tmp/nomad --bootstrap-expect=1
```

Nomad will also output some debugging information.

_We are running Nomad here as both a client and a server. In a production deployment, clients and servers are usually separate – the server is responsible for coordinating the cluster but doesn't actually run any jobs, while the client doesn't do any cluster management and just runs jobs. In development the distinction isn't important._

If everything worked as expected you should now be able to access UIs for Nomad and Consul on your local machine using these links:

[Consul UI](http://localhost:8500/)

[Nomad UI](http://localhost:4646/)

The Consul UI will show you all of the services which are currently available on the cluster – which for now, will just be Nomad and Consul itself! The Nomad UI will show you details of the Nomad cluster and any executing jobs.

But we don't have any jobs yet… so let's [move on to the next step and run a Nomad job!](./02-running-a-job.md)
