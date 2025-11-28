# go-rest-api-template
This is a template for creating golang REST API. Copy/Paste to use it. 

The main goal is to make it quick and easy to make a REST API
that has all the boilerplate sorted out, so you can go 
straight to customizing configuration and adding endpoints.

### Gen AI Friendly

This template includes Claud Code artifacts:

- context files for Go style guide
- context files for an Observability standard
- context files for Endpoint standards
- context files for openapi.yaml
- context files for reliability standards

### Dockerize

Template includes a Dockerfile.

## Configuration
Uses viper and cobra for handling configuration with
these features:

- environment variables
- command line args
- yaml config file

This template has an example of configuration
that would be typical of most projects. 

### Yaml config file

### Environment variables

### Command line args

## Makefile

This project includes a simple make file for compilation
and building the rest api binary.

To run it:

```bash
  echo "create a compile example "
  echo "create a binary example "
```

# Observability

The current template is logging-oriented for observability. 
We love logging because we control what observability events
and data we collect. It requires precision engineering, but
it rewards that in the end with the tightest observability that
can be exported to any tool, especially log management tools. 

This template includes directions for setting up log files
that are uploaded to Graphana's stack: Loki + Grafana dashboard. 

Also included is a sample Loki + Grafana dashboard that shows errors
and typical rest api stats.

## Logging
This project uses Go's structured logging `slog` library

The default for the template is to send log output to stdout and stderr. 
There is a configuration parameter that has these options:
- errors to  { stderr | stdout | logfile }
- info to { stderr | stdout | logfile }
- metrics to { stderr | stdout | logfile }
- traces to { stderr | stdout | logfile }

You can specify the log file path for any log file you give.
A new logfile is created each day to make it easier to prune log files. 
You can set log events to more than 1 place: (e.g. errors to stderr and logfile 'errs.jsonl')

Usually you will want your dev mode to log to files so you can read them.
Production typically uses stderr and stdout as targets. 

Note: A good future add-on could be an OTEL instrumented variant
with Grafana stack tracing and metrics. 

### Slog resources

- https://www.dash0.com/guides/logging-in-go-with-slog
- https://betterstack.com/community/guides/logging/logging-in-go/#final-thoughts





# Project Dependencies

Uses golang libraries:

Create dependencies
```bash
 go get github.com/jackc/pgx/v5/pgxpool
 go get -u github.com/spf13/cobra@latest
 go get github.com/spf13/viper
```