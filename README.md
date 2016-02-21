certs
=====

Certs is a common library for obtaining and managing TLS certificates with ACME CAs like Let's Encrypt. It is available for all platforms and has no external dependencies (not even libc!). Certs makes it possible to manage certificates in bulk and securely share them with other nodes in your infrastructure. Certs uses the excellent [acme package by xenolf](https://github.com/xenolf/lego).

It comes in two flavors: `certs`, a CLI tool, and `certsd`, a long-running daemon.

The default certificate authority is Let's Encrypt.

**This is a work in progress! The examples shown in this README are tentative and subject to change, as most features are not even implemented yet.**


## `certs`

Run `certs` for one-time commands, such as issuing, renewing, or revoking certificates. It exits once the task is complete.


#### Quick Start

Certs stores account information and other assets in `$HOME/.certs` (or `%HOMEPATH%\.certs` on Windows) by default.

First, register an account with the CA:

```
$ certs register -u email@example.com
```

To obtain certificates you have to prove ownership. You can currently do this 3 ways:

1. **http challenge:** This requires serving a resource over port 80.

2. **tls-sni challenge:** This requires serving a special TLS handshake over port 443.

3. **dns challenge:** This requires DNS credentials to set a temporary record in the zone file.

By default, certs will try http or tls-sni, so make sure those ports are available. If you provide DNS credentials through environment variables, however, certs will use the dns challenge. The dns challenge is nice because the domains don't have to be pointed at the machine you're running certs on.

Generate a certificate and key for example.com, which get bundled in PEM format to stdout:

```
$ certs issue -d example.com
```

Assuming that bundle was saved to a bundle.pem, you could renew the certificate:

```
$ cat bundle.pem | certs renew
```

If your private key is compromised, revoke the certificate thusly:

```
$ cat bundle.pem | certs revoke
```

To obtain a SAN certificate for example.com as common name and www.example.com as alternate:

```
$ certs issue -d example.com -d www.example.com
```

To obtain a certificate for which you already have a private key:

```
$ certs issue -d example.com -k key.pem
```

To obtain a certificate for which you already have a CSR:

```
$ certs issue -c csr.pem
```


## `certsd`

Run `certsd` to spawn a long-running child process that continuously keeps your certificates renewed. It also opens a port with an authenticated REST API so you can issue commands and securely transfer certificates and keys.

... todo ...
