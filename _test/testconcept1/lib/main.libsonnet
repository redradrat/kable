local kausal = import "ksonnet-util/kausal.libsonnet";

local deployment = kausal.apps.v1.deployment;
local container = kausal.core.v1.container;
local port = kausal.core.v1.containerPort;
local service = kausal.core.v1.service;

local grafanaDeploy(name) = deployment.new(
        name=name, replicas=2,
        containers=[
          container.new("grafana", "grafana/grafana")
          + container.withPorts([port.new("ui", 10330)]),
        ],
      );


// Final JSON Object
{
  new(name):: [
    grafanaDeploy(name),
    kausal.util.serviceFor(grafanaDeploy(name))
  ]
}
