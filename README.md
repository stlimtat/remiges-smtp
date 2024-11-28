# remiges-smtp
smtp client with file scrapping

## Tech stack
The following software were used to build this project and application:
- [Bazel](https://bazel.build)
  Bazel is used to build and run the application.
- [Cobra](https://github.com/spf13/cobra)
- [Viper](https://github.com/spf13/viper)
- [Zerolog](https://github.com/rs/zerolog)
- [mox](https://github.com/mjl-/mox)

## Usage
### Basic Usage
```
bazel run //cmd/smtpclient
```
### Lookup MX
Lookup MX DNS records for provided domain
```
bazel run //cmd/smtpclient lookupmx --domain=<domain with ending dot>
```
Example:
```
bazel run //cmd/smtpclient lookupmx --domain=dcs1.biz.
bazel run //cmd/smtpclient lookupmx --domain=remiges.tech.
```
