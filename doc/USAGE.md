# smtpclient cli
The following are the basic capabilities of the smtpclient cli


<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#basic-usage">Basic Usage</a></li>
    <ol>
        <li><a href="#lookup-mx">LookupMX</a></li>
        <li><a href="#gen-dkim">Gen DKIM</a></li>
    </ol>
  </ol>
</details>

## Basic Usage

We use bazel as the build engine/makefile mechanism for this project

```sh
bazel run //cmd/smtpclient
```

All commands (with the blessing of cobra) provide:
1. command line completion
1. help

```sh
bazel run //cmd/smtpclient help
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>

## Gen DKIM


## Lookup MX

Lookup MX DNS records for provided domain

```sh
bazel run //cmd/smtpclient lookupmx -- --domain=<domain with ending dot>
```

Example:

```sh
bazel run //cmd/smtpclient lookupmx -- --domain=dcs1.biz.
bazel run //cmd/smtpclient lookupmx -- --domain=remiges.tech.
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>
