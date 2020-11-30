# Magic Castle Plugins

This repository contains plugins used by [Magic Castle](https://www.github.com/computecanada/magic_castle) but
than cannot be simply deployed by [Puppet Magic Castle](https://www.github.com/computecanada/puppet-magic_castle)
because they require to be compiled first before being deployed.

List of plugins
- `cmd/consul2slurm`: Go plugin for consul-template that takes a JSON list of consul nodes config and
return Slurm node definition in [`slurm.conf`](https://slurm.schedmd.com/slurm.conf.html) format.
