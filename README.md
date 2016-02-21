certs
=====

Certs is a common library for obtaining and managing TLS certificates with ACME CAs like Let's Encrypt. It is available for all platforms and has no external dependencies (not even libc!). Certs makes it possible to manage certificates in bulk and securely share them with other nodes in your infrastructure. Certs uses the excellent [acme package by xenolf](https://github.com/xenolf/lego).

It comes in two flavors: `certs`, a CLI tool, and `certsd`, a long-running daemon.

The default certificate authority is Let's Encrypt.

**This is a work in progress! The examples shown in this README are tentative and subject to change, as most features are not even implemented yet.**


## `certs`

Run `certs` for one-time commands, such as issuing, renewing, or revoking certificates. It exits once the task is complete.


### Quick Start

Certs stores account information and other assets in `$HOME/.certs` (or `%HOMEPATH%\.certs` on Windows) by default.

First, register an account with the CA and agree to their terms:

```
$ certs register -u email@example.com --agree
```

To obtain certificates you have to prove ownership. You can currently do this 3 ways:

1. **http challenge:** This requires serving a resource over port 80.

2. **tls-sni challenge:** This requires serving a special TLS handshake over port 443.

3. **dns challenge:** This requires DNS credentials to set a temporary record in the zone file.

By default, certs will try http or tls-sni, so make sure those ports are available. If you provide DNS credentials through environment variables, however, certs will use the dns challenge. The dns challenge is nice because the domains don't have to be pointed at the machine you're running certs on and you don't need to give certs permission to bind to low ports. (TODO: Determine if certs will drop privileges; if not, recommend setcap!)

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

### Extended Tutorial

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

You can also obtain and renew certificates in bulk with certs, but since outputting a hundred+ certificates and just as many keys to stdout can be unwieldy, we write these to the file system instead. You can find the certificates and keys in `$HOME/.certs`. There, you'll see a lightweight database created that keeps track of where certificates are for which names, and the actual certificates and keys stored in subfolders. This folder structure and database is the same one used by `certsd`, so the tools are compatible.

To perform bulk issuance, you can pass in a CSV file, where each line is a certificate, and each entry on a line is a name to add to the certificate (first name is Common Name, the rest are Subject Alternative Names):

```
$ certs issue --csv "domains.csv"
```

When obtaining certificates in bulk, they're stored in the `$HOME/.certs` folder. If a domain fails to verify, the whole process exits with an error. Certificates for domains that already have a certificate will not be re-issued without the `-f` flag to force re-issuance. (TODO: Figure out precisely how we differentiate certificates -- whether by all SAN names or just CN...)

Or maybe you have a lot of certificates you need to obtain, but the challenge for each one has to be solved differently (maybe they're spread out across different DNS providers), you can load a JSON file that gives you total control over each certificate to issue:

```
$ certs issue --json "domains.json"
```

(TODO: Example JSON structure.)

Renewal in bulk is the same, except run `certs renew` instead of `certs issue`. When renewing, only domains that are within 30 days of expiration will be renewed. You can adjust this window with the `--days` option.



## `certsd`

Run `certsd` to spawn a long-running child process that continuously keeps your certificates renewed. It also opens a port with an authenticated REST API so you can issue commands and securely transfer certificates and keys.

Note that if certsd is running and using ports 80 and 443 (which it does by default), you can still use certs to solve http and tls-sni challenges.

Certsd implements privilege de-escalation, so you can safely run as root to bind low ports, and it will immediately drop privileges in the child process.

### Quick start

First you'll need to register an account with your CA if you haven't already, and agree to their terms. You can do this with certs:

```
$ certs register -u email@example.com --agree
```

To start certsd with default settings, just run it, either as root or with setcap so it can bind to low ports:

```
$ certsd
```

This opens ports 80 and 443 for serving ACME challenges. It also will maintain certificates in the `$HOME/.certs` folder, keeping them renewed. Errors will be logged to stderr (use the --log option to change this).

(TODO: A certs flag to not keep the domain renewed. Sets a flag in its entry in the database.)

Certsd also provides an API for managing certificates. To enable the API, start certsd with a file called `certsd.conf` in the working directory, or use the `--conf` option to load it from somewhere else. The config file must authorize at least one user, since all API requests must be authenticated.

(TODO: Example YAML file)

(TODO: API examples)

(TODO: USR1 to reload config)
