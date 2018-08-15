# 04 - Configuration management

Now that we've got a job running in Nomad, and a service registered in Consul, let's look at how we can use Consul to manage our application's configuration so they can talk to each other.

## Creating another job

Our new "hello world" service is getting a bit big, so we really need to look at splitting it into microservices so we can scale it up. Let's build a web frontend that responds to user requests, and use our existing job as an internal service we can use to find out what text to respond with. 

Let's make a new job spec in `hello-world-frontend.hcl`:

```hcl
job "hello-world-frontend" {
  datacenters = ["dc1"]
  group "web" {
    task "server" {
      driver = "docker"
      config {
        image = "altmetric/cakefindr-hello-world-frontend"
      }
      resources {
        network {
          port "http" {}
        }
      }
      template {
        destination = "local/env"
        env         = true
        data = <<EOH
PORT={{ env "NOMAD_PORT_http" }}
BACKENDS={{range service "hello-world" }}{{ .Address }}:{{ .Port }},{{ end }}
EOH
      }
      service {
        name = "hello-world-frontend"
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

```

This is a bit different from our first job file, because it contains a `template` stanza instead of an `env` stanza. This allows us to use [Consul Template](https://github.com/hashicorp/consul-template)… uh, templates as part of our Nomad jobs specification. Let's have a closer look at that stanza with some additional comments:

```
// Tell Nomad to generate a file inside the task using this template
template {

  // Specify the location that this file will be generated at in the environment
  // of this task
  destination = "local/env"

  // Tell Nomad to use this file to define environment variables for the task
  env = true

  // data contains a heredoc that holds the template data. Embedded Consul 
  // Template directives use the {{ double curly }} syntax and support a rich
  // language for formatting text.
  //
  // Our template extracts the port that Nomad has selected for this task and
  // puts it into the `PORT` environment variable. We then ask Consul for a list
  // of all available `hello-world` services that we'll be using as our backends
  // and joins them with commas into the `BACKENDS` environment variable, which
  // is used by the frontend service to locate backends.
  data = <<EOH
PORT={{ env "NOMAD_PORT_http" }}
BACKENDS={{range service "hello-world" }}{{ .Address }}:{{ .Port }},{{ end }}
EOH
}
```

This is a similar approach to the one we currently use when build `.pam_environment` files for our existing infrastructure.

Run your new frontend server using Nomad…

```bash
→ nomad run ./hello-world-frontend.hcl
```

… and you should be able to resolve the service using Consul:

```
→ dig srv +noall +answer hello-world-frontend.service.consul @127.0.0.1 -p 8600
hello-world-frontend.service.consul. 0 IN SRV 1 1 20172 0a000015.addr.dc1.consul.
```

We can then access our new front-end service using this address:

```
→ curl 0a000015.addr.dc1.consul.:20172
Hello world!
```

Now we're webscale.

_In a real deployment system, we would include a dedicated load-balancer like haproxy to distribute requests between Nomad allocations, instead of including the backend configuration inside the application. That's beyond the scope of this exercise, because we're still trying to figure out exacly how it will work!_

## Additional configuration

It's possible that users will want to receive some different text instead of the boring old "Hello world!" message. Fortunately, our backend service can take a `RESPONSE_TEXT` environment variable to customise its response.

But we don't want to have to deploy a whole new job every time we want to change the text! Instead can use Consul's key-value store to make changes on the fly and allow Nomad to reload our application when needed.

Let's start by creating a configuration variable in Consul using the `consukl kv` command:

```
→ consul kv put config/response_text "Goodbye cruel world"
Success! Data written to: config/response_text
```

Then we can update our backend application to start using this configuration data. Remove the `env` stanza from `hello-world.hcl` and replace it with a Consul Template configuration similar to the frontend:

```
template {
  data = <<EOH
PORT={{ env "NOMAD_PORT_http" }}
RESPONSE_TEXT="{{ key "config/response_text" }}"
EOH
  destination = "local/env"
  env         = true
}
```

We use the Consul Template `key` directive to insert a value from the KV store.

We can now resubmit the backend job to Nomad:

```bash
→ nomad run ./hello-world.hcl
```

Then make a request to our frontend again. Note that the port may change, as Nomad will relaunch the application configured to talk to backends with the new configuration we have supplied:

```
→ dig srv +noall +answer hello-world-frontend.service.consul @127.0.0.1 -p 8600
hello-world-frontend.service.consul. 0 IN SRV 1 1 20212 0a000015.addr.dc1.consul.

→ curl 0a000015.addr.dc1.consul.:20212
Goodbye cruel world
```

Success! Now we have a fully dynamic, scalable, configurable microservice-driven hello world application!

Nomad will apply changes from Consul in realtime, so we can update the KV store and see changes without having to relaunch our applications:

```
→ consul kv put config/response_text "Hello, nice to meet you"
Success! Data written to: config/response_text

→ dig srv +noall +answer hello-world-frontend.service.consul @127.0.0.1 -p 8600
hello-world-frontend.service.consul. 0 IN SRV 1 1 28037 0a000015.addr.dc1.consul.

→ curl 0a000015.addr.dc1.consul.:28037
Hello, nice to meet you
```

_You might find that your services are unavailable for a short period while Nomad restarts containers with updated environment variables. In a real deployment, we include configuration which makes Nomad implement a rolling restart approach to ensure that applications remain available when updates to configuration are being deployed._

_It might be worth trying to see if you can include some information in the `RESPONSE_TEXT` variable that demonstrates that requests to the frontend are being fulfilled by different backend servers. Try updating the template configuration for the backends to include a reference to the `NOMAD_ALLOC_ID` environment variable!_

## Next steps

Now we've figured out how to run applications on our local cluster and have them talk to each other. Now let's [build a larger cluster!](./05-joining-a-cluster.md)
