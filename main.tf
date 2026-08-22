terraform {
  required_version = "1.15.8"

  cloud {

    organization = "kevin-labs"

    workspaces {
      name = "cli-testing"
    }
  }
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.23.0"
    }
  }
}

provider "google" {
  project     = "kevin-lol-service-staging"
  region      = "northamerica-northeast2"
  zone        = "northamerica-northeast2-b"
}

resource "google_sql_database_instance" "main" {
  name             = "main-instance"
  database_version = "POSTGRES_18"

  settings {
    tier              = "db-f1-micro"
    edition           = "ENTERPRISE"
    availability_type = "ZONAL"
    disk_type         = "PD_HDD"
    disk_size         = 10
    disk_autoresize   = false

    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = false
    }
  }
}
