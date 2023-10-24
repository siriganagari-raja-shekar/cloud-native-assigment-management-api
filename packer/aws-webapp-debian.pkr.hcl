packer {
  required_plugins {
    amazon = {
      version = ">= 0.0.2"
      source  = "github.com/hashicorp/amazon"
    }
  }
}

variable "aws_region" {
  type = string
}
variable "dev_subnet_id" {
  type = string
}
variable "source_ami_id" {
  type = string
}
variable "ssh_username" {
  type = string
}

variable "linux_group" {
  type = string
}
variable "linux_user" {
  type = string
}
variable "user_home_dir" {
  type = string
}
variable "standard_log_file" {
  type = string
}
variable "error_log_file" {
  type = string
}
variable "account_csv_path" {
  type = string
}


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
      "STANDARD_LOG_FILE=${var.standard_log_file}",
      "ERROR_LOG_FILE=${var.error_log_file}",
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1"
    ]
  }

}
