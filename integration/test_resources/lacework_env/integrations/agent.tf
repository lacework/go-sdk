resource "random_id" "instance_id" {
  byte_length = 4
}

data "aws_ami" "ubuntu1804" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"]
}

resource "aws_key_pair" "key" {
  key_name   = var.aws_key_pair_name
  public_key = file(pathexpand("~/.ssh/id_rsa.pub"))
}

resource "aws_instance" "ubuntu1804" {
  connection {
    host        = aws_instance.ubuntu1804.public_ip
    user        = "ubuntu"
    private_key = file(var.aws_key_pair_file)
  }

  ami                         = data.aws_ami.ubuntu1804.id
  instance_type               = "t2.medium"
  key_name                    = aws_key_pair.key.key_name
  subnet_id                   = var.subnet_id
  vpc_security_group_ids      = ["${var.security_group_id}"]
  associate_public_ip_address = true

  tags = {
    Name       = "go_sdk_cli_test_resource_${random_id.instance_id.hex}"
    X-Dept     = "Engineering"
    X-Customer = "tech-ally"
    X-Project  = "go-sdk/lacework-cli"
    X-Contact  = "tech-ally@lacework.net"
    X-TTL      = "730"
  }

  provisioner "remote-exec" {
    inline = [
      "sudo hostname ubuntu1804-lacework",
      "sudo apt-get install -y curl",
      "sudo curl -o /tmp/install.sh ${var.install_sh_url}",
      "sudo chmod +x /tmp/install.sh",
      "sudo /tmp/install.sh ${lacework_agent_access_token.token.token}"
    ]
  }
}
