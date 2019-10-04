package config

var defaultConfig = `
cache:
  expiration: 1h

discovery:
  directory: ""

docker:
  registries: {}

log:
  level: info

http:
  address: :8080

prometheus:
  path: /metrics

ui:
  enabled: true
  staticPath: ui/static
  templatePath: ui/template.html.mustache
`
