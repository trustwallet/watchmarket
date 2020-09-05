# -*- mode: Python -*-

# For more on Extensions, see: https://docs.tilt.dev/extensions.html
load('ext://restart_process', 'docker_build_with_restart')

# Records the current time, then kicks off a server update.
# Normally, you would let Tilt do deploys automatically, but this
# shows you how to set up a custom workflow that measures it.
compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api ./cmd/api/api && go build -o build/worker'
if os.name == 'nt':
  compile_cmd = 'build.bat'

local_resource(
  'ci',
  compile_cmd,
  deps=['./cmd'])

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

k8s_yaml('tilt-deployments/deployment.yaml')
k8s_resource('nginx-proxy', port_forwards=8081,
             resource_deps=['ci'])
