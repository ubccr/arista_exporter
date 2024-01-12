# Arista Exporter

Prometheus exporter for Arista switches. Uses the Go Client for eAPI to query
output from various commands related to MLAG status, supervisor redundancy
status, port-channel status, and power supply status. This exporter is intended
to query multiple Arista switches from an external host.

The `/arista` metrics endpoint exposes the Arista metrics and requires a
`target` parameter.  The `module` parameter can also be used to select which
probe commands to run, the default module is `power`. Available modules are:

- mlag
- power
- portchannel
- redundancy
- switchover

The `/metrics` endpoint exposes Go and process metrics for this exporter.

## Configuration

This exporter requires an eapi.conf file. More details [see
here](https://github.com/aristanetworks/goeapi#getting-started). Example
config:

```ini
[connection:arista1.example.com]
host=arista1.example.com
username=admin
password=root
enablepwd=passwd
transport=https

[connection:arista2.example.com]
host=arista2.example.com
username=admin
password=root
enablepwd=passwd
transport=https
```

## Prometheus configs

```yaml
- job_name: arista
  metrics_path: /arista
  static_configs:
  - targets:
    - arista1.example.com
    - arista2.example.com
    labels:
      module: mlag,power,portchannel
  - targets:
    - arista3.example.com
    labels:
      module: redundancy,switchover,power,portchannel
  relabel_configs:
  - source_labels: [__address__]
    target_label: __param_target
  - source_labels: [__param_target]
    target_label: instance
  - source_labels: [module]
    target_label: __param_module
  - target_label: __address__
    replacement: 127.0.0.1:9465
```

Example systemd unit file [here](systemd/arista_exporter.service)

## Sample Metrics

```
arista_power_supply_state{powerSupply="1",state="ok"} 1
arista_power_supply_state{powerSupply="2",state="ok"} 1
arista_redundancy_all_agent_sso_ready 1
arista_redundancy_communication_desc{status="Up"} 1
arista_redundancy_last_mode_change_time{reason="Supervisor has control of the active supervisor lock"} 1.6874010372967267e+09
arista_redundancy_mode{status="active"} 1
arista_redundancy_peer_mode{status="standby"} 1
arista_redundancy_slot_id{unitDesc="Primary"} 1
arista_redundancy_switchover_count 0
arista_redundancy_switchover_ready 1
arista_mlag_config_sanity{status="consistent"} 1
arista_mlag_detail{mlagState="secondary",peerMlagState="primary"} 1
arista_mlag_failover{cause=""} 0
arista_mlag_last_state_change 336.392112271
arista_mlag_local_inf_status{status="up"} 1
arista_mlag_neg_status{status="connected"} 1
arista_mlag_peer_link{status="up"} 1
arista_mlag_ports{state="activeFull"} 2
arista_mlag_ports{state="activePartial"} 0
arista_mlag_ports{state="configured"} 0
arista_mlag_ports{state="disabled"} 0
arista_mlag_ports{state="inactive"} 8
arista_mlag_state{state="active"} 1
arista_mlag_state_changes 2
arista_portchannel_ports{interface="Port-Channel1",state="active"} 0
arista_portchannel_ports{interface="Port-Channel1",state="inactive"} 2
arista_portchannel_ports{interface="Port-Channel2",state="active"} 1
arista_portchannel_ports{interface="Port-Channel2",state="inactive"} 1
arista_portchannel_ports{interface="Port-Channel3",state="active"} 2
arista_portchannel_ports{interface="Port-Channel3",state="inactive"} 0
```

## License

arista_exporter is released under the Apache License Version 2.0. See the LICENSE file.
