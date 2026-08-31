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
  project = var.project
  region  = var.region
  zone    = var.zone
}

resource "google_sql_database_instance" "default" {
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

resource "google_artifact_registry_repository" "default" {
  repository_id = "kevin-lol-service"
  description   = "Docker repository for ${local.service_name}"
  format        = "DOCKER"
}

resource "google_cloud_run_v2_service" "lol_service"  {
  location = var.region
  name     = local.service_name

  template {
    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.default.connection_name]
      }
    }

    containers {
      image = "${google_artifact_registry_repository.default.location}-docker.pkg.dev/${var.project}/${google_artifact_registry_repository.default.repository_id}/${local.image}"
      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }

      env {
        name = "KEVIN_RIOT_API_KEY"
      }
    }
  }
}

resource "google_cloud_run_v2_job" "migrate" {
  location = var.region
  name     = "${local.service_name}-db-migrate"

  template {
    template {
      volumes {
        name = "cloudsql"
        cloud_sql_instance {
          instances = [google_sql_database_instance.default.connection_name]
        }
      }

      containers {
        image = "${google_artifact_registry_repository.default.location}-docker.pkg.dev/${var.project}/${google_artifact_registry_repository.default.repository_id}/${local.image}"
        args = ["migrate"]

        env {
          name = "KEVIN_RIOT_API_KEY"
          value_from {
            secret_version = google_secret_manager_secret.riot_api_key.secret_id
          }
        }
      }
    }
  }
}

resource "google_secret_manager_secret" "riot_api_key" {
  secret_id = "riot-api-key"

  replication {
    auto {}
  }
}

resource "google_secret_manager_secret_version" "secret-version-data" {
  secret = google_secret_manager_secret.secret.name
  secret_data = "secret-data"
}

resource "google_secret_manager_secret_iam_member" "secret_access" {
  secret_id = google_secret_manager_secret.riot_api_key.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${data.google_project.project.number}-compute@developer.gserviceaccount.com"

  depends_on = [
    google_secret_manager_secret.riot_api_key
  ]
}