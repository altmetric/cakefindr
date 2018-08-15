# 06 – The Cake Is Not A Lie – Team Alpha (Gio & Shane)

Your team's job is to run the application in the Docker image `altmetric/cakefindr-keyserver-3`. This is an application that can provide part of the secret key required to decrypt the location of the cake in the UI.

Use the techniques we learned earlier to run the Docker image on the Nomad cluster, and register a service called `keyserver-3`. The image is configured using two environment variables: `PORT` contains the port on which the application listens (like the examples we tried earlier) and `SECRET_KEY` contains the secret key that should be returned by the application when fetching from the path `/key-3`. The value returned by the application should be `KEY_HERE`. Good luck!
