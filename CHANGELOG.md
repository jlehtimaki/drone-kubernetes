# Changelog

All notable changes will be shown here.

Versioning follows Semantic Versioning principals.

## v.1.3.0 - 2020-10-15
### Added
- Support for rollout status checks. This is ON by default.
Users can disable this feature by adding setting `rollout_check: false`
Users can change timeout for rollout check with `rollout_timeout`

## v1.2.2 - 2020-09-16
### Bug fix
- Fixed diff in kustomize command

## v1.2.0 - 2020-09-01
### Changes
- Dockerfile changed to support multiarch builds using `docker buildx`
- Rewrote downloadFile function
- Rewrote pipelining kustomize commands
- Rewrote `chmod` command to use `os.Chmod`

## v1.1.0 - 2020-08-13
### Added
- Added diff as new command
- Redid the commands part of the code


## v1.0.0 - 2020-04-24
Renamed the whole repository and added modularity so that this can be used in different Kubernetes deployments
### Added
- Users can set the type of Kubernetes deployment EKS/Baremetal
- New environment variables: PLUGIN_TYPE, PLUGIN_CA, PLUGIN_TOKEN, PLUGIN_K8S_USER and PLUGIN_K8S_SERVER\
These enable the Baremetal functionality to the plugin
