# zerotier-consumer-go-challenge
Coding Challenge for ZeroTier (Broccoli IT) Webhook Consumer written in GoLang

1. I worked on this locally using IntelliJ Ultimate on Pengwin WSL2 (WLinux)
2. The run was installed using my Google Cloud Platform account

# GoLang on Pengwin WSL2
* download [go1.22.5.linux-amd64.tar.gz](https://go.dev/dl/go1.22.5.linux-amd64.tar.gz)
```shell
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xvf <downloadFolder>/go1.22.5.linux-amd64.tar.gz
export GOROOT=/usr/local/go
go version go1.22.5 linux/amd64
```

# Docker on Pengwin WSL2
```shell
sudo apt install ca-certificates curl gnupg lsb-release
curl -fsSL https://download.docker.com/linux/debian/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
chmod -v 700 /home/matthew/.gnupg
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" \
  | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt update
sudo apt upgrade
sudo apt install docker-ce docker-ce-cli containerd.io
sudo usermod -aG docker $USER
sudo /etc/init.d/docker start
# re-login for the Unix Group change
```

# Docker Build
```shell
# build
docker build .
docker images # find image ID
# test
docker run -p 4444:4444 a4b0dae69173
```

