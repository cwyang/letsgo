#!/bin/bash
set -eu

TIMEZONE=Asia/Seoul
USERNAME=greenlight

read -p "Enter password for greenlight DB user: " DB_PASSWORD

export LC_ALL=ko_KR.UTF-8

## script logic

add-apt-repository --yes universe
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all

useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

rsync --archive --chown="${USERNAME}:${USERNAME}" /root/.ssh /home/${USERNAME}

ufw allow 22
ufw allow 80/tcp
ufw allow 443/tcp
ufw --force enable

apt --yes install fail2ban

#install migrate CLI tool
#curl -L https://...

apt --yes install postgresql

sudo -i -u postgres psql -c "CREATE DATABASE greenlight"
sudo -i -u postgres psql -d greenlight -c "CREATE extension if not exists citext"
sudo -i -u postgres psql -d greenlight -c "CREATE role greenlight with login password '${DB_PASSWORD}'"

echo "GREENLIGHT_DB_DSN='postgres://greenlight:${DB_PASSWORD}@localhost/greenlight'" >> /etc/environment

#install caddy
#...

echo "Script complete. rebooting!"
reboot
