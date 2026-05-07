resource "yandex_vpc_network" "lab2_network" {
  name = "lab2-network"
}

resource "yandex_vpc_subnet" "lab2_subnet" {
  name           = "lab2-subnet"
  zone           = "ru-central1-b"
  network_id     = yandex_vpc_network.lab2_network.id
  v4_cidr_blocks = ["10.0.1.0/24"]
}

resource "yandex_compute_instance" "lab2_vm" {
  name        = "lab2-vm"
  platform_id = "standard-v3"
  resources {
    cores  = 2
    memory = 4
  }

  boot_disk {
    initialize_params {
      # Ubuntu 22.04 LTS (актуальный ID нужно проверить в консоли)
      image_id = "fd80on0ma1ic60hees6n"
    }
  }

  network_interface {
    subnet_id = yandex_vpc_subnet.lab2_subnet.id
    nat       = true
  }

  metadata = {
    ssh-keys = "ubuntu:${file("~/.ssh/id_rsa.pub")}"
  }
}