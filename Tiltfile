docker_build('moficodes/k8s-admin', '.')

k8s_yaml('k8s/app.yaml')

k8s_resource('admin-deploy', port_forwards='9000')