---
outline: deep
---
# Запуск с использованием docker-compose

В этом разделе показано, как установить услугу Nutrix с помощью [docker-compose](https://docs.docker.com/compose/).

::: tip
На Windows убедитесь, что [WSL2](https://learn.microsoft.com/en-us/windows/wsl/install) установлен перед началом.
:::

::: info
Убедитесь, что установлены [git](https://git-scm.com/), [docker](https://www.docker.com/) и [docker-compose](https://docs.docker.com/compose/)

:::

### 1. Клонируйте репо pos и запустите docker compose
```sh
git clone https://github.com/nutrixpos/pos.git
cd pos
docker-compose up --build -d
```

### 2. Настройка автозапуска

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

### 3. Откройте консоль nutrix
Откройте консоль nutrix по адресу http://localhost:3000