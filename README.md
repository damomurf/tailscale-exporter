tailscale-exporter
==================

Prometheus Metrics and Metadata via the (Beta) Tailscale API.

Requires a Tailscale API Token available from the Management Web UI under Keys -> API Key -> "Generate One-off key"

Execute the exporter as follows:

```
$ ./tailscale-exporter -tailnet <tailnet> -token <api-token>
```

This will listen by default on port 8080.

Currently the exporter generates the following metrics for each device in your Tailscale network:

```
tailscale_blocks_incoming{id="123456",name="hostname.domain"} 1              
tailscale_device_info{external="false",hostname="hostname",id="123456",name="hostname.domain"} 1
tailscale_expires{id="123456",name="hostname.domain"} 1.620480219e+09        
tailscale_external{id="123456",name="hostname.domain"} 1
tailscale_last_seen{id="123456",name="hostname.domain"} 1.614856637e+09      
tailscale_upgrade_available{id="123456",name="hostname.domain"} 1
```
