output "ansible_inventory" {
  description = "Готовый inventory-файл для Ansible"
  value = <<EOF
[lab_vms]
lab2_vm ansible_host=${yandex_compute_instance.lab2_vm.network_interface[0].nat_ip_address} ansible_user=ubuntu ansible_ssh_private_key_file=~/.ssh/id_rsa
EOF
}
