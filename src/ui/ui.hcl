job "ui" {
  datacenters = ["dc1"]
  group "web" {
    task "server" {
      driver = "docker"
      config {
        image = "altmetric/cakefindr-ui"
      }
      resources {
        network {
          port "http" {}
        }
      }
      template {
        data = <<EOH
PORT={{ env "NOMAD_PORT_http" }}
KEYSERVER_1="{{range service "keyserver-1" }}{{ .Address }}:{{ .Port }}{{ end }}"
KEYSERVER_2="{{range service "keyserver-2" }}{{ .Address }}:{{ .Port }}{{ end }}"
KEYSERVER_3="{{range service "keyserver-3" }}{{ .Address }}:{{ .Port }}{{ end }}"
EOH
        destination = "local/env"
        env         = true
      }
      service {
        name = "ui"
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
