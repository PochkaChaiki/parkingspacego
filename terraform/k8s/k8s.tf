resource "yandex_vpc_network" "lab3_network" {
  name = "lab3_network"
}

resource "yandex_vpc_subnet" "lab3_subnet" {
  name           = "lab3-subnet"
  zone           = "ru-central1-b"
  network_id     = yandex_vpc_network.lab3_network.id
  v4_cidr_blocks = ["10.0.1.0/24"]
}

resource "yandex_kubernetes_cluster" "lab3_k8s" {
  name        = "lab3-k8s-cluster"
  description = "Kubernetes cluster for parking app"
  
  network_id = yandex_vpc_network.lab3_network.id
  
  release_channel = "REGULAR"

  master {
    zonal {
      zone      = "ru-central1-b"
      subnet_id = yandex_vpc_subnet.lab3_subnet.id
    }
    
    public_ip = true
    
    maintenance_policy {
      auto_upgrade = true
    }
  }

  service_account_id      = yandex_iam_service_account.k8s_sa.id
  node_service_account_id = yandex_iam_service_account.k8s_node_sa.id
}

resource "yandex_kubernetes_node_group" "lab3_nodes" {
  name        = "lab3-node-group"
  cluster_id  = yandex_kubernetes_cluster.lab3_k8s.id
  
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
      subnet_ids = [yandex_vpc_subnet.lab3_subnet.id]
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
