# -*- mode: Python -*-

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')

# Records the current time, then kicks off a server update.
# Normally, you would let Tilt do deploys automatically, but this
# shows you how to set up a custom workflow that measures it.
# compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api ./cmd/api/api && go build -o build/worker'
# if os.name == 'nt':
#   compile_cmd = 'build.bat'

# local_resource(
#   'ci',
#   compile_cmd,
#   deps=['./cmd'])

# go build -o ./bin/api ./cmd/api/main.go
# go build -o ./bin/worker ./cmd/worker/main.go
# docker build -t trust/watchmarket:seed-local -f seed/Dockerfile seed/
# docker build -t trust/watchmarket:proxy-local -f nginx/Dockerfile nginx/
# docker build -t trust/watchmarket:pg-health-local -f scripts/pg-check/Dockerfile scripts/pg-check/

# building go binary
build_api = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/api ./cmd/api/main.go'
build_worker = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/worker ./cmd/worker/main.go'

local_resource(
  'api-build',
  build_api,
  deps=['./cmd']
)

local_resource(
  'worker-build',
  build_worker,
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
k8s_resource('nginx-proxy', port_forwards=8081)
