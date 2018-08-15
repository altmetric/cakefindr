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
        data = <<EOH
PORT={{ env "NOMAD_PORT_http" }}
BACKENDS={{range service "hello-world" }}{{ .Address }}:{{ .Port }},{{ end }}
EOH
        destination = "local/env"
        env         = true
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
