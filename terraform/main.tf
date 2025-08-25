provider "google" {
  credentials = "../gcp-sa-key.json"
  project     = var.project_id
  region      = var.region_name
}

terraform {
  backend "gcs" {
    credentials = "../gcp-sa-key.json"
  }

  required_version = "= 1.9.0"
}

resource "google_artifact_registry_repository" "repo" {
  location      = var.region_name
  repository_id = var.docker_image_name
  format        = "DOCKER"
}
