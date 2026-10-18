---
outline: deep
---
# التشغيل باستخدام docker-compose

يوضح هذا القسم كيفية تثبيت خدمة Nutrix باستخدام [docker-compose](https://docs.docker.com/compose/).

::: tip
على Windows تأكد من تثبيت [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install) قبل المتابعة.
:::

::: info
تأكد من تثبيت [git](https://git-scm.com/) و [docker](https://www.docker.com/) و [docker-compose](https://docs.docker.com/compose/)

:::

### 1. استنسخ مستودع pos وقم بتشغيل docker compose
```sh
git clone https://github.com/nutrixpos/pos.git
cd pos
docker-compose up --build -d
```

### 2. تكوين البدء التلقائي

::: info
For **Windows** users !

Enable systemd in wsl config

```sh
vim /etc/wsl.conf
```
```/etc/wsl.conf
[boot]
systemd=true
```

You will then need to close the Linux distribution. You can use the command `wsl.exe --shutdown` in PowerShell to restart all WSL instances.

Once your Linux distribution restarts, systemd will be running. You can confirm using the command: systemctl list-unit-files --type=service, which will show the status of any services associated with your Linux distribution.

for more info refer to [WSL Systemd](https://learn.microsoft.com/en-us/windows/wsl/systemd)
:::

Create systemd unit files
```sh
sudo nano /etc/systemd/system/nutrix.service
```

``` nutrix.service{9}
[Unit]
Description=My Docker Compose Service
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/path/to/your/pos_repo
ExecStartPre=/usr/bin/docker-compose down
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
Restart=no

[Install]
WantedBy=multi-user.target
```
::: tip
Make sure you change **/path/to/your/pos_repo** in the above file to the pos repo path on your machine
:::

Enable and start the service

```sh
sudo systemctl enable nutrix.service
sudo systemctl start nutrix.service
```

### 3. افتح وحدة تحكم nutrix
افتح وحدة تحكم nutrix على http://localhost:3000