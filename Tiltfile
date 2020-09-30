# -*- mode: Python -*-

local_resource(
  'lint+tests',
  'make go-lint && make go-test && go-integration',
  trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False
)

local_resource(
  'seed',
  'docker build -t trust/watchmarket:seed-local -f deployment/utils/seed/Dockerfile . && (kubectl delete -f deployment/utils/seed/job.yaml || echo 1) && kubectl apply -f deployment/utils/seed/job.yaml && echo Start seeding',
  trigger_mode=TRIGGER_MODE_MANUAL, auto_init=False
)
 
docker_build("trust/watchmarket:api-local", ".", build_args={"SERVICE":"api"})
docker_build("trust/watchmarket:worker-local", ".", build_args={"SERVICE":"worker"})

yaml = helm(
  'deployment/charts/watchmarket',
  name='local',
  namespace='tilt-watchmarket-local',
  values=['./deployment/charts/watchmarket/values.local.yaml']
)

# k8s namespace bootstrap
local('kubectl create namespace tilt-watchmarket-local || echo 1')

k8s_yaml(yaml)
k8s_resource('api', port_forwards=8421)

k8s_resource('postgres', port_forwards=8585)
k8s_resource('redis', port_forwards=8586)
