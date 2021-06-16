The `squirrel.service` file is an example systemd unit file for keeping the squirrel service running.


- change the configuration flags under `[Service]`.  
- `sudo cp squirrel.service /etc/systemd/system/`
- `sudo systemctl start squirrel.service`

# Useful commands
Reload Unit File:  
```
sudo systemctl daemon-reload
```


Start/Stop/Restart/status:

```
sudo systemctl start squirrel.service
sudo systemctl stop squirrel.service
sudo systemctl restart squirrel.service
sudo systemctl status squirrel.service
```
# Useful variables
Environment="HTTP_PROXY=http://proxy.company.com:port/"
Environment="HTTPS_PROXY=http://proxy.company.com:port/"

Tail logs:
```
journalctl -u squirrel.service -f
```
