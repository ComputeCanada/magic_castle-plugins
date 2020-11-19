#!/usr/bin/python
import sys
import json

slurmds = json.loads(sys.argv[1])
if slurmds:
  configset = sorted(set((c['ServiceMeta']['gpus'], c['ServiceMeta']['realmemory'], c['ServiceMeta']['cpus']) for c in slurmds))
  weight_map = {k:weight for weight, k in enumerate(configset, 1)}

  for c in slurmds:
    print "NodeName={} CPUs={} RealMemory={} Gres=gpu:{} Weight={}".format(
        c['Node'],
        c['ServiceMeta']['cpus'],
        c['ServiceMeta']['realmemory'],
        c['ServiceMeta']['gpus'],
        weight_map[(c['ServiceMeta']['gpus'], c['ServiceMeta']['realmemory'], c['ServiceMeta']['cpus'])]
    )
