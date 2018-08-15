# 06 – The Cake Is Not A Lie – Team Four (Paul & Maciej)

Your team's job is to run the application in the Docker image `altmetric/cakefindr-ui`. This is a simple UI that decrypts an internal secret using data sources from three different keyservers.

Use the techniques we learned earlier to run the Docker image on the Nomad cluster, and register a service called `ui`. The image is configured using four environment variables:

- `PORT` contains the port on which the application listens (like the examples we tried earlier)
- `KEYSERVER_1` contains an IP address and port pair where the application can access secret key number 1 (e.g. `10.0.0.1:7623`). The key will be downloaded from the endpoint `/key-1`, and this service will be available as Consul service `keyserver-1`.
- `KEYSERVER_2` contains an IP address and port pair where the application can access secret key number 1 (e.g. `10.0.0.1:7623`). The key will be downloaded from the endpoint `/key-2`, and this service will be available as Consul service `keyserver-2`.
- `KEYSERVER_3` contains an IP address and port pair where the application can access secret key number 1 (e.g. `10.0.0.1:7623`). The key will be downloaded from the endpoint `/key-3`, and this service will be available as Consul service `keyserver-3`.

When this application runs, you should be able to access the service at `ui.service.consul` and decrypt the secret. Good luck!
