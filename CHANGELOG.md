# Changelog

All notable changes will be shown here.

Versioning follows Semantic Versioning principals.

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
