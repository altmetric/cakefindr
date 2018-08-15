# 03 - Service discovery

Now we've got a running Nomad cluster hosting our highly scalable "hello world" service. We need to figure out how to make this rich API accessible to consumers in an environment where ports and machines might change unexpectedly – so let's use Consul.

## Interacting with Consul

Consul implements service discovery using its own protocol where applications can subscribe to changes, as an HTTP API, and as a DNS server. We already looked briefly at the [Consul UI](http://localhost:8500), but the primary usage of Consul is through direct integration or its use as a DNS server.

The basic Consul CLI is limited, but we can use it to see what services are currently registered:

```
→ consul catalog services
consul
nomad
nomad-client
```

We can use the DNS API to query Consul for more details about the location of services. Every service registered in Consul receives a DNS name with the pattern `<service-name>.service.consul`. We currently only have Nomad and Consul itself registered, so let's look at those by using `dig` to speak directly to the Consul DNS server on port 8600:

```
→ dig +noall +answer nomad.service.consul @127.0.0.1 -p 8600
nomad.service.consul. 0 IN  A 10.0.0.21

→ dig +noall +answer consul.service.consul @127.0.0.1 -p 8600
consul.service.consul.  0 IN  A 10.0.0.21
```

Okay, that's not too surprising. Consul and Nomad are both running on your local machine, so they resolve to `127.0.0.1`. But we don't know what port they're on!

Fortunately, Consul supports [DNS SRV records](https://en.wikipedia.org/wiki/SRV_record) – these are special DNS records that can include additional information about resolved names. Importantly, this includes the port that they're running on!

We can ask `dig` to retrieve `SRV` records using the `srv` flag:

```
→ dig srv +noall +answer nomad.service.consul @127.0.0.1 -p 8600
nomad.service.consul. 0 IN  SRV 1 1 4646 0a000015.addr.dc1.consul.
nomad.service.consul. 0 IN  SRV 1 1 4647 0a000015.addr.dc1.consul.
nomad.service.consul. 0 IN  SRV 1 1 4648 0a000015.addr.dc1.consul.
```

Okay, that's already looking more useful. We can see that the Nomad service in Consul has three ports registered for the different network services it provides. The HTTP one is probably most useful for us, so we can ask Consul for those details:

```
→ dig srv +noall +answer http.nomad.service.consul @127.0.0.1 -p 8600
http.nomad.service.consul. 0  IN  SRV 1 1 4646 0a000015.addr.dc1.consul.
```

Even better. Now we know where to access the HTTP endpoint for Nomad.

Let's make it a bit easier to use by configuring the machine to use Consul for all name resolution when a domain ends in `.consul`!

On a Mac, we can run this command:

```
sudo mkdir -p /etc/resolver && sudo tee /etc/resolver/consul > /dev/null <<EOF
nameserver 127.0.0.1
port 8600
EOF
```

> Requirements on Linux might be different depending on your distribution!

We should immediately be able to ping services using their Consul name:

```
→ ping http.nomad.service.consul
PING http.nomad.service.consul (10.0.0.21): 56 data bytes
64 bytes from 10.0.0.21: icmp_seq=0 ttl=64 time=0.057 ms
```

Now we can access services that are registered in Consul. But how do we make our own services appear?

## Registering a service

Nomad integrates tightly with Consul to allow applications to discover each other. To register a Nomad job in Consul, we can add a `service` stanza to the job file we generated earlier. Add the following stanza inside the `task` block in your job file:

```
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
```

This block tells Nomad to register a Consul service for our job, using the `http` port label that we already defined. It also includes a "service check", which Consul can use to check if the application is working correctly.

Let's submit the updated job to Nomad:

```
→ nomad run ./hello-world.hcl
```

The job will be updated in-place, and we'll immediately see the new service in Consul:

```
→ consul catalog services
consul
hello-world
nomad
nomad-client
```

We can also query the Consul DNS API for details:

```
→ dig srv +noall +answer hello-world.service.consul @127.0.0.1 -p 8600
hello-world.service.consul. 0 IN  SRV 1 1 23326 0a000015.addr.dc1.consul.
hello-world.service.consul. 0 IN  SRV 1 1 26036 0a000015.addr.dc1.consul.
hello-world.service.consul. 0 IN  SRV 1 1 26706 0a000015.addr.dc1.consul.
hello-world.service.consul. 0 IN  SRV 1 1 27903 0a000015.addr.dc1.consul.
hello-world.service.consul. 0 IN  SRV 1 1 21984 0a000015.addr.dc1.consul.
```

So now we know where to find all of our allocations, just by using the DNS API. They'll also be visible in the [Consul UI](http://localhost:8500/ui/dc1/services/hello-world).

_Try changing the `count` parameter in the Nomad job file and re-running the job, then running the above DNS query again. Nomad will keep Consul up-to-date with the real allocations that are running._

If a health check fails, Consul will immediately remove the failing task from the responses it gives to queries. This provides an effective method of load-balancing, and makes sure that clients don't get routed to dead instances of services.

## Next steps

Now we can create jobs in Nomad, and register them with Consul so that other applications can find them. But how do we actually use this information in other applications? Let's look at [using Consul for configuration management](./04-configuration-management.md).
