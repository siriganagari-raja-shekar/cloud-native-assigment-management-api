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
variable "postgres_host" {}
variable "postgres_port" {}
variable "postgres_user" {}
variable "postgres_password" {}
variable "postgres_db" {}
variable "account_csv_path" {}
variable "server_port" {}

source "amazon-ebs" "webapp" {
  ami_name      = "csye6225-webapp-debian-aws"
  instance_type = "t2.micro"
  region        = "${var.aws_region}"
  source_ami_filter {
    filters = {
      image-id            = "ami-06db4d78cb1d3bbf9"
      root-device-type    = "ebs"
      virtualization-type = "hvm"
    }
    most_recent = true
    owners      = ["136693071363"]
  }

  ami_users    = ["149723291571"]
  subnet_id    = "${var.dev_subnet_id}"
  ssh_username = "admin"
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
      "POSTGRES_HOST=${var.postgres_host}",
      "POSTGRES_PORT=${var.postgres_port}",
      "POSTGRES_USER=${var.postgres_user}",
      "POSTGRES_PASSWORD=${var.postgres_password}",
      "POSTGRES_DB=${var.postgres_db}",
      "ACCOUNT_CSV_PATH=${var.account_csv_path}",
      "SERVER_PORT=${var.server_port}"
    ]
  }

}
