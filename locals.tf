locals {
  service_name = "lol-service"

  image_tag = "latest"
  image = "${local.service_name}:${local.image_tag}"
}
