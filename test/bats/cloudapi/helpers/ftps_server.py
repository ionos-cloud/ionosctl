#!/usr/bin/env python3
"""Minimal FTPS server for ionosctl upload integration tests.

Starts a local FTP server with Explicit TLS (AUTH TLS) on a given port,
using a self-signed certificate. Writes its PID to a file for teardown.

Usage:
    python3 ftps_server.py <port> <root_dir> <certfile> <keyfile> <pidfile> <user> <password>

The server creates iso-images/ and hdd-images/ directories under root_dir
automatically.
"""

import os
import sys

from pyftpdlib.authorizers import DummyAuthorizer
from pyftpdlib.handlers import TLS_FTPHandler
from pyftpdlib.servers import FTPServer


def main():
    if len(sys.argv) != 8:
        print(f"Usage: {sys.argv[0]} <port> <root_dir> <certfile> <keyfile> <pidfile> <user> <password>")
        sys.exit(1)

    port = int(sys.argv[1])
    root_dir = sys.argv[2]
    certfile = sys.argv[3]
    keyfile = sys.argv[4]
    pidfile = sys.argv[5]
    user = sys.argv[6]
    password = sys.argv[7]

    # Create upload directories
    for d in ("iso-images", "hdd-images"):
        os.makedirs(os.path.join(root_dir, d), exist_ok=True)

    authorizer = DummyAuthorizer()
    authorizer.add_user(user, password, root_dir, perm="elradfmw")

    handler = TLS_FTPHandler
    handler.authorizer = authorizer
    handler.certfile = certfile
    handler.keyfile = keyfile
    handler.tls_control_required = False  # Allow AUTH TLS (Explicit TLS)
    handler.tls_data_required = False
    handler.passive_ports = range(60000, 60100)

    server = FTPServer(("127.0.0.1", port), handler)
    server.max_cons = 50

    # Write PID so tests can kill us
    with open(pidfile, "w") as f:
        f.write(str(os.getpid()))

    print(f"FTPS server listening on 127.0.0.1:{port}, root={root_dir}")
    server.serve_forever()


if __name__ == "__main__":
    main()
