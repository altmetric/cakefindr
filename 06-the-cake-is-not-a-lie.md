# 06 – The Cake Is Not A Lie 🍰

Now we've learned how to build a Nomad and Consul cluster, run and configure jobs, and join a larger cluster. It's time to implement a distributed application that will help us find some cake.

Somewhere in the office is a hidden cake that we can all enjoy if we sucessfully deploy `cakefindr` to our Nomad cluster. Since this is obviously a high-security application, the location of the cake is stored in an encrypted format, and to decrypt it we'll need to load keys from three different servers into a front-end application.

## The applications

### Team One (Ana & Anna)

Team One will be running the application `altmetric/cakefindr-keyserver-1`. This application needs to serve one of the secret keys used to decode the secret location.

[More details](./06-the-cake-is-not-a-lie-team-one.md)

### Team Two (Jonathan & Måel)

Team Two will be running the application `altmetric/cakefindr-keyserver-2`. This application needs to serve one of the secret keys used to decode the secret location.

[More details](./06-the-cake-is-not-a-lie-team-two.md)

### Team Three (Gio & Shane)

Team Three will be running the application `altmetric/cakefindr-keyserver-3`. This application needs to serve one of the secret keys used to decode the secret location.

[More details](./06-the-cake-is-not-a-lie-team-three.md)

### Team Four (Paul & Maciej)

Team Four will be running the application `altmetric/cakefindr-ui`. This application needs to access data from the three keyservers to decode it's insternal payload containing the secret location.

[More details](./06-the-cake-is-not-a-lie-team-four.md)

## The Cake 🍰

Follow the details in the linked pages above – and when we all sucessfully deploy our applications to our shared Nomad cluster, Team Aleph should be able to decode and display the secret location!
