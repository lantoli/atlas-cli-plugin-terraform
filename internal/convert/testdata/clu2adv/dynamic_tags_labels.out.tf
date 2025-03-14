resource "mongodbatlas_advanced_cluster" "simplified" {
  project_id   = var.project_id
  name         = "cluster"
  cluster_type = "REPLICASET"
  replication_specs = [
    {
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          priority      = 7
          electable_specs = {
            node_count    = 3
            instance_size = "M10"
          }
        }
      ]
    }
  ]
  tags = var.tags

  # Generated by atlas-cli-plugin-terraform.
  # Please confirm that all references to this resource are updated.
}

resource "mongodbatlas_advanced_cluster" "expression" {
  project_id   = var.project_id
  name         = "cluster"
  cluster_type = "REPLICASET"
  replication_specs = [
    {
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          priority      = 7
          electable_specs = {
            node_count    = 3
            instance_size = "M10"
          }
        }
      ]
    }
  ]
  tags = {
    for key, value in local.tags : key => replace(value, "/", "_")
  }

  # Generated by atlas-cli-plugin-terraform.
  # Please confirm that all references to this resource are updated.
}

resource "mongodbatlas_advanced_cluster" "expression_individual" {
  project_id   = var.project_id
  name         = "cluster"
  cluster_type = "REPLICASET"
  replication_specs = [
    {
      region_configs = [
        {
          provider_name = "AWS"
          region_name   = "US_EAST_1"
          priority      = 7
          electable_specs = {
            node_count    = 3
            instance_size = "M10"
          }
        }
      ]
    }
  ]
  tags = merge(
    {
      for key, value in var.tags : key => replace(value, "/", "_")
    },
    {
      tag1    = var.tag1val
      "tag 2" = var.tag2val
    }
  )

  # Generated by atlas-cli-plugin-terraform.
  # Please confirm that all references to this resource are updated.
}
