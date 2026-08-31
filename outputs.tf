output "lol-service-endpoint" {
  value = google_cloud_run_v2_service.lol_service.uri
}
