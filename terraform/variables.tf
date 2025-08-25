variable "docker_image_name" {
  description = "Name of the Docker image"
  type        = string
}

variable "project_id" {
  description = "Google Cloud Project ID"
  type        = string
}

variable "region_name" {
  description = "Region for the Cloud Run service"
  type        = string
}