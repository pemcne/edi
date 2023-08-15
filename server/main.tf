provider "google" {
  project = var.gcp_project
  region = "us-central1"
  zone = "us-central1-c"
}

resource "google_compute_instance" "vm" {
  name = "bots"
  machine_type = "e2-micro"
  zone = "us-central1-c"
  boot_disk {
    initialize_params {
      image = "ubuntu-2204-lts"
    }
  }
  network_interface {
    network = "default"
    access_config {}
  }
  service_account {
    email = "slack-bot@sinuous-aviary-391419.iam.gserviceaccount.com"
    scopes = [
      "https://www.googleapis.com/auth/cloud-platform"
    ]
  }
  timeouts {}
  metadata = {
    "startup-script" = file("startup.sh")
  }
  lifecycle {
    ignore_changes = [ network_interface ]
  }
}

variable "gcp_project" {
  type = string
}