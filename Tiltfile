# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')

go_compile_cmd = 'sh tilt-build/go-build.sh'
if os.name == 'nt':
  go_compile_cmd = 'tilt-build\\go-build.bat'

local_resource(
  'go-build',
  go_compile_cmd,
  deps=['./cmd','./api','./config','./db','./internal','./pkg','./redis','./services']
)

local_resource(
  'ci-test',
  'make go-test & make go-integration',
  trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False
)

docker_build("trust/watchmarket:seed-local", "seed", dockerfile="seed/Dockerfile")
docker_build("trust/watchmarket:proxy-local", "nginx", dockerfile="nginx/Dockerfile")
docker_build("trust/watchmarket:pg-health-local", "scripts/pg-check", dockerfile="scripts/pg-check/Dockerfile")

docker_build_with_restart(
  'trust/watchmarket:api-local',
  '.',
  build_args={"SERVICE":"linux/api"},
  entrypoint=["/app/main", "-c", "/config/config.yml"],
  dockerfile='Dockerfile.runner',
  only=[
    './bin/linux/api','./config.yml',
  ],
  live_update=[
    sync('./bin/linux/api','/app/main')
  ]
)

docker_build_with_restart(
  'trust/watchmarket:worker-local',
  '.',
  build_args={"SERVICE":"linux/worker"},
  entrypoint=["/app/main", "-c", "/config/config.yml"],
  dockerfile='Dockerfile.runner',
  only=[
    './bin/linux/worker','./config.yml',
  ],
  live_update=[
    sync('./bin/linux/worker','/app/main')
  ]
)

yaml = helm(
  'charts/watchmarket',
  # The release name, equivalent to helm --name
  name='local',
  # The namespace to install in, equivalent to helm --namespace
  namespace='tilt-watchmarket-local',
  # The values file to substitute into the chart.
  values=['./charts/watchmarket/values.local.yaml']
  )
  
# k8s namespace bootstrap
k8s_namespace_create = 'kubectl create namespace tilt-watchmarket-local || true'
if os.name == 'nt':
  k8s_namespace_create = 'kubectl create namespace tilt-watchmarket-local || echo 1'
local(k8s_namespace_create)

k8s_yaml(yaml)
k8s_resource('nginx-proxy', port_forwards=8081, 
             resource_deps=['go-build'])

k8s_resource('postgres', port_forwards=8585)

k8s_resource('redis', port_forwards=8586)
