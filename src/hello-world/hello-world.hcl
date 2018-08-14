// This is a simple "hello world" Nomad job with some information about what
// each stanza does.

// The top level of every job file is the "job" stanza. This defines the name
// of the job as it will appear in Nomad, Consul, and to other applications.
// This has to be a DNS-friendly name, so you can use letter, numnbers, and
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

    count = 5

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

      service {
        name = "hello-world"
        port = "http"
        check {
          type     = "http"
          path     = "/ping"
          interval = "10s"
          timeout  = "2s"
        }
      }
    }
  }
}
