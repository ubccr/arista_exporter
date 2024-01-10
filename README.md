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
  - target_label: __address__
    replacement: 127.0.0.1:9465
  - source_labels: [module]
    target_label: __param_module
```

Example systemd unit file [here](systemd/arista_exporter.service)

## License

arista_exporter is released under the Apache License Version 2.0. See the LICENSE file.
