job "keyserver-2" {
  datacenters = ["dc1"]
  group "web" {
    task "server" {
      driver = "docker"
      config {
        image = "altmetric/cakefindr-keyserver-2"
      }
      resources {
        network {
          port "http" {}
        }
      }
      template {
        data = <<EOH
      PORT={{ env "NOMAD_PORT_http" }}
      SECRET_KEY="{{ key "config/keyserver-2/secret" }}"
      EOH
        destination = "local/env"
        env         = true
      }

      service {
        name = "keyserver-2"
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
