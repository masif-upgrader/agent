## About

The **Masif Upgrader agent** is a component of *Masif Upgrader*.

Consult Masif Upgrader's [manual] on its purpose
and the agent's role in its architecture
and [demo] for a full stack live demonstration.

## Configuration

The configuration file (usually `/etc/masif-upgrader-agent/config.ini`)
looks like this:

```ini
[interval]
check=300
report=300
retry=60

[master]
host=infra-mgmt.intern.example.com:8150

[tls]
cert=/var/lib/puppet/ssl/certs/mail.example.com.pem
key=/var/lib/puppet/ssl/private_keys/mail.example.com.pem
ca=/var/lib/puppet/ssl/certs/ca.pem

[log]
level=info
```

The *interval* section defines several intervals:

 option | description
 -------|---------------------------------------------------------------------------------------------------------------------------
 check  | Check every x seconds whether any packages can be upgraded
 report | Once any packages can be upgraded, report the set of required actions to upgrade all of them every x seconds to the master
 retry  | If any action fails, retry it after x seconds (0 or not set = don't retry anything)

*master.host* is the master's address (HOST:PORT).

The *tls* section describes the X.509 PKI:

 option | description
 -------|---------------------------------------------------
 cert   | TLS client certificate chain (may include root CA)
 key    | TLS client private key
 ca     | TLS server root CA certificate

*log.level* defines the logging verbosity and is one of:

* error
* warning
* info
* debug

[manual]: https://github.com/masif-upgrader/manual
[demo]: https://github.com/masif-upgrader/demo
