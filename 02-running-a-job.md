# 02 - Running a job

So we've got a Nomad cluster running locally on your machine. But it's not that useful yet, because it's not running anything!

## Creating a job file

To start running applications on the cluster, we have to give Nomad a "job file" written according to the [job specification format](https://www.nomadproject.io/docs/job-specification/index.html). This describes various attributes of the job – like which application we want to run, how much memory needs to be allocated to it, and how many instances of the job need to be executed.

Job files are written in the [Hashicorp Configuration Language (HCL)](https://github.com/hashicorp/hcl), which is a JSON-like format used by all of the Hashicorp tools. Let's start by running a basic "hello world" job to get the feel for how this works.

Create a file called `hello-world.hcl` with the following content, and have a look at the content – it hopefully explains what the different parts of the file are for.

```hcl
// This is a simple "hello world" Nomad job with some information about what
// each stanza does.

// The top level of every job file is the "job" stanza. This defines the name
// of the job as it will appear in Nomad, Consul, and to other applications.
// This has to be a DNS-friendly name, so you can use letter, numbers, and
// hyphens only.
job "hello-world" {

  // We define the datacenter in which this job will run. For the moment, we
  // only have a single datacenter, but Nomad supports job specifications that
  // enforce things like ensuring that jobs run across multiple datacenters for
  // redundancy.
  datacenters = ["dc1"]

  // Each job contains one or more "groups". A group defines a collection of
  // tasks that need to be co-located on the same server. Most of the time this
  // will contain a single task.
  group "web" {

    // The "task" stanza defines an individual unit of work that will be
    // executed for this job. In this case, we are running a simple Docker image
    // containing a web server.
    task "server" {

      // Tell Nomad that we are creating a job using the Docker driver.
      driver = "docker"

      // Configuration is specific to each driver. In this case, we are telling
      // the Docker driver to use the specified Docker image.
      config {
        image = "altmetric/cakefindr-hello-world"
      }

      // For each task, we define the resources which are required to run it.
      // This can include things like CPU, memory, and network bandwidth. In
      // this case, we tell Nomad that this task needs a network port to be
      // allocated to it, and that we will refer to this port with the label
      // "http" (there's no magic to this value).
      resources {
        network {
          port "http" {}
        }
      }

      // The Docker image we are running requires a HTTP port to be passed in an
      // environment variable called PORT. We can use the "env" stanza to define
      // this environment variable – Nomad allows us to use string interpolation
      // to access the port number that it has chosen to use for this task.
      env {
        "PORT" = "${NOMAD_PORT_http}"
      }
    }
  }
}
```


## Planning the job

Once we have created the job file, we can tell Nomad to start running the job on your cluster. The first step is to ask Nomad to plan execution of the job – this will check that the specification is valid, and the cluster has the resources to run the job.

We can use the `nomad` command-line utility to interact with the Nomad cluster. Plan your job using the `nomad plan` command:

```bash
→ nomad plan ./hello-world.hcl
```

Nomad will respond with an execution plan which describes the actions it's going to take:

```
+ Job: "hello-world"
+ Task Group: "web" (1 create)
  + Task: "server" (forces create)

Scheduler dry-run:
- All tasks successfully allocated.

Job Modify Index: 0
To submit the job with version verification run:

nomad job run -check-index 0 ./hello-world.hcl

When running the job with the check-index flag, the job will only be run if the
server side version matches the job modify index returned. If the index has
changed, another user has modified the job and the plan's results are
potentially invalid.
```

From the output, you can see that Nomad has detected that it needs to create a new instance of the web server task we defined in the job file, because the job isn't currently running. The CLI includes a job versioning feature that can help to avoid the risk of overwriting other users' jobs, but we don't have to worry about that for now.

## Running the job

When the job has been sucessfully planned, we can ask Nomad to actually start running it:

```
→ nomad run ./hello-world.hcl
==> Monitoring evaluation "6caf8763"
    Evaluation triggered by job "hello-world"
    Allocation "c8685372" created: node "830325ef", group "web"
    Evaluation status changed: "pending" -> "complete"
==> Evaluation "6caf8763" finished with status "complete"
```

Nomad will start an `evaluation` process, which is the term it uses to describe the creation or update process for a job file. In our case, this will result in the creation of an `allocation`, which is the concrete instantiation of the tasks that we included in our job.

We can check that the job started successfully using the `nomad status` command:

```
→ nomad status
ID           Type     Priority  Status   Submit Date
hello-world  service  50        running  2018-08-14T22:44:56+01:00
```

## Examining the job

We can find out more information about an individual job using `nomad status <job-name>`:

```
→ nomad status hello-world
ID            = hello-world
Name          = hello-world
Submit Date   = 2018-08-14T22:44:56+01:00
Type          = service
Priority      = 50
Datacenters   = dc1
Status        = running
Periodic      = false
Parameterized = false

Summary
Task Group  Queued  Starting  Running  Failed  Complete  Lost
web         0       0         1        0       0         0

Allocations
ID        Node ID   Task Group  Version  Desired  Status   Created    Modified
c8685372  830325ef  web         0        run      running  4m33s ago  4m21s ago
```

This command outputs a full status report for the job we submitted, including a breakdown of all the individual `task`s and `allocation`s that it contains. We can see that Nomad is running a single instance of the `web` group that we specified, and a single `allocation` to provide the server.

We can examine each individual allocation to find out more information using `nomad status <allocation-id>`:

```
→ nomad status c8685372
ID                  = c8685372
Eval ID             = 6caf8763
Name                = hello-world.web[0]
Node ID             = 830325ef
Job ID              = hello-world
Job Version         = 0
Client Status       = running
Client Description  = <none>
Desired Status      = run
Desired Description = <none>
Created             = 8m59s ago
Modified            = 8m47s ago

Task "server" is "running"
Task Resources
CPU        Memory           Disk     IOPS  Addresses
0/100 MHz  824 KiB/300 MiB  300 MiB  0     http: 10.0.0.21:31721

Task Events:
Started At     = 2018-08-14T21:44:59Z
Finished At    = N/A
Total Restarts = 0
Last Restart   = N/A

Recent Events:
Time                       Type        Description
2018-08-14T22:44:59+01:00  Started     Task started by client
2018-08-14T22:44:56+01:00  Driver      Downloading image altmetric/cakefindr-hello-world:latest
2018-08-14T22:44:56+01:00  Task Setup  Building Task Directory
2018-08-14T22:44:56+01:00  Received    Task received by client
```

This command describes the individual allocation in detail, including the details of the resources it is using. Have a look at the output of your own `nomad status` command, and you'll see the `Task Resources` section tells us the port that Nomad has selected for this application. We can use this information to access the application that we're running:

```
→ curl 10.0.0.21:31721
Hello world!
```

## Modifying the job

When we want to make changes to a job – for example, using a new Docker image or changing the number of instances we deploy – we can update the job file and re-submit it to Nomad.

Let's run more servers so that our hello world service is properly webscale.

Add a new line to your job file immediately inside the `group` stanza:

```
count = 5
```

This will tell Nomad to ensure it always runs five instances of our server task. We can now re-plan this job with Nomad and see what happens!

```
→ nomad plan ./hello-world.hcl
+/- Job: "hello-world"
+/- Task Group: "web" (4 create, 1 in-place update)
  +/- Count: "1" => "5" (forces create)
      Task: "server"

Scheduler dry-run:
- All tasks successfully allocated.
```

We can see from the output that Nomad wants to create four new instances of our web server because the `count` has changed. Let's go ahead and apply that change:

```
→ nomad run ./hello-world.hcl
==> Monitoring evaluation "e53227f8"
    Evaluation triggered by job "hello-world"
    Allocation "3d5abd33" created: node "830325ef", group "web"
    Allocation "61db9b69" created: node "830325ef", group "web"
    Allocation "6276bbca" created: node "830325ef", group "web"
    Allocation "e614f937" created: node "830325ef", group "web"
    Allocation "c8685372" modified: node "830325ef", group "web"
    Evaluation status changed: "pending" -> "complete"
==> Evaluation "e53227f8" finished with status "complete"
```

Now we can check that we're running more instances of this server:

```
→ nomad status hello-world
…

Summary
Task Group  Queued  Starting  Running  Failed  Complete  Lost
web         0       0         5        0       0         0

Allocations
ID        Node ID   Task Group  Version  Desired  Status   Created    Modified
34d43449  830325ef  web         0        run      running  3s ago     0s ago
5ea4f06c  830325ef  web         0        run      running  3s ago     0s ago
caea11d6  830325ef  web         0        run      running  3s ago     1s ago
ee0fdfe8  830325ef  web         0        run      running  3s ago     0s ago
c8685372  830325ef  web         0        run      running  6m1s ago   6m1s ago
```

We can see that we now have _five_ allocations running, including the one we already had running.

_Go ahead and look at the output of `nomad status <allocation-id>` for the different allocations. You'll see that they've each been allocated their own port by Nomad, and that all of these ports will offer the same HTTP endpoint._

## Stopping the job

When we're done with a job, we can tell Nomad to shut it down with `nomad stop <job-name>`

```bash
→ nomad stop hello-world
```

Nomad will stop all running `allocation`s and end the job. We will still be able to see it in the output of `nomad status` and `nomad status hello-world`:

```]
→ nomad status
ID           Type     Priority  Status          Submit Date
hello-world  service  50        dead (stopped)  2018-08-14T23:19:10+01:00

→ nomad status hello-world
…
Summary
Task Group  Queued  Starting  Running  Failed  Complete  Lost
web         0       0         0        0       5         0
…
```

Nomad marks the job as complete, and periodically cleans up the list of dead jobs for us. We can start the job again by re-submitting the job file:

```bash
→ nomad run ./hello-world.hcl
```

Nomad will create new allocations for this job and start it up again.

_Now might be a good opportunity to poke around the [Nomad UI](http://localhost:4646). You can see a useful representation of the [job](http://localhost:4646/ui/jobs/hello-world) and details of all the [allocations](http://localhost:4646/ui/jobs/hello-world/web). This might be easier than using the CLI in some cases._

## Next steps

So now we've covered the steps required to run the Hashicorp stack and deploy applications to it. But we've also seen that Nomad will assign a random port number to our job, and if we had a larger cluster it might choose to run it on a different machine. How can we find out programatically where the instances of each job are running?

Well, we can use Nomad's integration with Consul to register services that an application provides, and make them accessible to other applications. Let's look next at [service discovery with Consul](./03-service-discovery.md).
