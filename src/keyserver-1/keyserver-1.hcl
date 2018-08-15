job "keyserver-1" {
  datacenters = ["dc1"]
  group "web" {
    task "server" {
      driver = "docker"
      config {
        image = "altmetric/cakefindr-keyserver-1"
      }
      resources {
        network {
          port "http" {}
        }
      }
      template {
        data = <<EOH
      PORT={{ env "NOMAD_PORT_http" }}
      SECRET_KEY="{{ key "config/keyserver-1/secret" }}"
      EOH
        destination = "local/env"
        env         = true
      }

      service {
        name = "keyserver-1"
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
