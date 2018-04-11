#!/bin/bash

# Script to run {.beat_name} in foreground with the same path settings that
# the init script / systemd unit file would do.

/usr/share/smshubbeat/bin/smshubbeat \
  -path.home /usr/share/smshubbeat \
  -path.config /etc/smshubbeat \
  -path.data /var/lib/smshubbeat \
  -path.logs /var/log/smshubbeat \
  $@

