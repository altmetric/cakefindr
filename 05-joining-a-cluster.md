# 05 – Joining a cluster

We've learned how to run Nomad and Consul, and deploy jobs on our platform. But one of the biggest benefits of this stack is being able to run multiple cluster nodes that can tolorate failures – so now we're going to build a cluster together!

*Make sure you're on the `Digital Science` wifi network before starting this, or it won't work!*

## Cleaning up

First, we'll stop the existing Consul and Nomad agents and clean up thier configuration. Close the running applicationa and remove their data directories:

```bash
rm -rf /tmp/consul /tmp/nomad
```

## Starting Consul

Now we'll start the Consul server again, but this time we'll include a `retry-join` parameter that tells it about an existing member of the cluster. Consul will use this address to join the cluster.

```bash
consul agent --data-dir=/tmp/consul -server -ui -retry-join=10.40.50.241
```

Consul will start and 🤞🏻 join our shared cluster. We can check that this worked using the `consul members` command to show all members in the cluster:

```
→ consul members
Node             Address            Status  Type    Build  Protocol  DC   Segment
Machine-1.local  10.216.3.206:8301  alive   server  1.2.2  2         dc1  <all>
Machine-2.local  10.216.3.217:8301  alive   server  1.2.2  2         dc1  <all>
```

## Starting Nomad

Now that you have a functioning Consul cluster member, we can start Nomad as well. Nomad will automatically use Consul to discover other Nomads in the cluster, so we don't need to supply any additional information:

```bash
nomad agent -client -server -data-dir=/tmp/nomad
```

We can check to confirm that our local Nomad server has joined the cluster:

```
→ nomad server members
Name                    Address       Port  Status  Leader  Protocol  Build  Datacenter  Region
Machine-1.local.global  10.216.3.206  4648  alive   true    2         0.8.4  dc1         global
Machine-2.local.global  10.216.3.217  4648  alive   false   2         0.8.4  dc1         global
```

And we can also check which clients are now available to run applications on using the `nomad node status` command:


```
→ nomad node status
ID        DC   Name            Class   Drain  Eligibility  Status
b489fa27  dc1 Machine-1.local  <none>  false  eligible     ready
a0bf5707  dc1 Machine-2.local  <none>  false  eligible     ready
```

## Next steps

Now we have a functioning Hashicorp cluster, let's [work together to deploy a distributed application!](./06-the-cake-is-not-a-lie.md)
