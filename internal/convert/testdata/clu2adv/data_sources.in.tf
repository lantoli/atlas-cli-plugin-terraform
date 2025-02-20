data "mongodbatlas_cluster" "singular" {
  # data source content is kept - singular
  project_id = mongodbatlas_advanced_cluster.example.project_id
  name       = mongodbatlas_advanced_cluster.example.name
  depends_on = [mongodbatlas_privatelink_endpoint_service.example_endpoint]
}

data "mongodbatlas_clusters" "plural" {
  # data source content is kept - plural
  project_id = mongodbatlas_advanced_cluster.example.project_id
}

data "mongodbatlas_advanced_cluster" "adv_singular" {
  # adv_cluster is not changed - adv_singular
  project_id = mongodbatlas_advanced_cluster.example.project_id
  name       = mongodbatlas_advanced_cluster.example.name
}

data "mongodbatlas_advanced_cluster" "adv_plural" {
  # adv_cluster is not changed - adv_plural
  project_id = mongodbatlas_advanced_cluster.example.project_id
}

resource "random_resource" "random1" {
  # other resources are left unchanged - random1
  hi1 = "there1"
}

data "random_datasource" "random2" {
  # other resources are left unchanged - random2
  hi2 = "there2"
}

# comments out of resources are kept

unknown_block "hello" {
  # unkown block types are kept
}

unknown_block {
  # plugin doesn't panic with unlabeled blocks - unknown_block
}

resource {
  # plugin doesn't panic with unlabeled blocks - resource
}

data {
  # plugin doesn't panic with unlabeled blocks - data
}
