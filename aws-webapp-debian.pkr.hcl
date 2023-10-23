packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "aws_region" {}
variable "dev_subnet_id" {}
variable "source_ami_id" {}
variable "ssh_username" {}

variable "linux_group" {}
variable "linux_user" {}
variable "user_home_dir" {}
variable "account_csv_path" {}


source "amazon-ebs" "webapp" {
  ami_name      = "csye6225-webapp-debian-${formatdate("YYYY_MM_DD_hh_mm_ss", timestamp())}"
  instance_type = "t2.micro"
  region        = "${var.aws_region}"
  source_ami_filter {
    filters = {
      image-id            = "${var.source_ami_id}"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["136693071363"]
  }

  ami_users    = ["149723291571"]
  subnet_id    = "${var.dev_subnet_id}"
  ssh_username = "${var.ssh_username}"
}

build {

  name = "build-webapp-ami"
  sources = [
    "source.amazon-ebs.webapp"
  ]

  provisioner "file" {
    source      = "app.zip"
    destination = "/tmp/app.zip"
  }

  provisioner "shell" {
    scripts = ["install.sh"]
    environment_vars = [
      "LINUX_GROUP=${var.linux_group}",
      "LINUX_USER=${var.linux_user}",
      "USER_HOME_DIR=${var.user_home_dir}",
      "ACCOUNT_CSV_PATH=${var.account_csv_path}",
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1"
    ]
  }

}
