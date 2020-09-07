# -*- mode: Python -*-

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')

# building go binary

local_resource(
  'api-build',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/api ./cmd/api/main.go',
  deps=['./cmd']
)

local_resource(
  'worker-build',
  'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/worker ./cmd/worker/main.go',
  deps=['./cmd']
)

docker_build("trust/watchmarket:seed-local", "seed", dockerfile="seed/Dockerfile")
docker_build("trust/watchmarket:proxy-local", "nginx", dockerfile="nginx/Dockerfile")
docker_build("trust/watchmarket:pg-health-local", "scripts/pg-check", dockerfile="scripts/pg-check/Dockerfile")

docker_build_with_restart(
  'trust/watchmarket:api-local',
  '.',
  build_args={"SERVICE":"api"},
  entrypoint=["/app/main", "-c", "/config/config.yml"],
  dockerfile='Dockerfile.runner',
  only=[
    './bin/','./config.yml',
  ],
  live_update=[
    sync('./bin/api','/app/main'),
  ],
)

docker_build_with_restart(
  'trust/watchmarket:worker-local',
  '.',
  build_args={"SERVICE":"worker"},
  entrypoint=["/app/main", "-c", "/config/config.yml"],
  dockerfile='Dockerfile.runner',
  only=[
    './bin/','./config.yml',
  ],
  live_update=[
    sync('./bin/worker','/app/main'),
  ],
)

yaml = helm(
  'charts/watchmarket',
  # The release name, equivalent to helm --name
  name='local',
  # The namespace to install in, equivalent to helm --namespace
  namespace='default',
  # The values file to substitute into the chart.
  values=['./charts/watchmarket/values.local.yaml']
  )

k8s_yaml(yaml)
k8s_resource('nginx-proxy', port_forwards=8081, 
             resource_deps=['api-build', 'worker-build'])
