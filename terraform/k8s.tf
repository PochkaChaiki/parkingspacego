resource "yandex_kubernetes_cluster" "lab2_k8s" {
  name        = "lab2-k8s-cluster"
  description = "Kubernetes cluster for parking app"
  
  network_id = yandex_vpc_network.lab2_network.id
  
  master {
    regional {
      region = "ru-central1"
      location {
        zone      = "ru-central1-b"
        subnet_id = yandex_vpc_subnet.lab2_subnet.id
      }
    }
    
    version = "1.28"
    public_ip = true
    
    maintenance_policy {
      auto_upgrade = true
    }
  }

  service_account_id      = yandex_iam_service_account.k8s_sa.id
  node_service_account_id = yandex_iam_service_account.k8s_node_sa.id
}

resource "yandex_kubernetes_node_group" "lab2_nodes" {
  name        = "lab2-node-group"
  cluster_id  = yandex_kubernetes_cluster.lab2_k8s.id
  version     = "1.28"
  
  instance_template {
    platform_id = "standard-v3"
    
    resources {
      memory = 4
      cores  = 2
    }
    
    boot_disk {
      type = "network-ssd"
      size = 30
    }
    
    network_interface {
      subnet_ids = [yandex_vpc_subnet.lab2_subnet.id]
      nat        = true
    }
    
    container_runtime {
      type = "containerd"
    }
  }
  
  scale_policy {
    auto_scale {
      min     = 1
      max     = 5
      initial = 2
    }
  }
  
  allocation_policy {
    location {
      zone = "ru-central1-b"
    }
  }
  
  maintenance_policy {
    auto_upgrade = true
    auto_repair  = true
  }
}

resource "yandex_iam_service_account" "k8s_sa" {
  name = "k8s-cluster-sa"
}

resource "yandex_iam_service_account" "k8s_node_sa" {
  name = "k8s-node-sa"
}

resource "yandex_resourcemanager_folder_iam_member" "k8s_sa_editor" {
  folder_id = var.folder_id
  role      = "editor"
  member    = "serviceAccount:${yandex_iam_service_account.k8s_sa.id}"
}

resource "yandex_resourcemanager_folder_iam_member" "k8s_node_viewer" {
  folder_id = var.folder_id
  role      = "viewer"
  member    = "serviceAccount:${yandex_iam_service_account.k8s_node_sa.id}"
}
