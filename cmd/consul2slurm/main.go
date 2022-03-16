package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Node structure
type Node struct {
	Name    string    `json:"Node"`
	Address string    `json:"NodeAddress"`
	Specs   NodeSpecs `json:"ServiceMeta"`
	Prefix  string
	Index   int
}

// NodeSpecs structure
type NodeSpecs struct {
	Cpus         int `json:"cpus,string"`
	Gpus         int `json:"gpus,string"`
	RealMemory   int `json:"realmemory,string"`
	MemSpecLimit int `json:"memspeclimit,string"`
	Weight       int `json:"weight,string"`
}

func main() {
	arg := []byte(os.Args[1])

	var nodes []Node

	if err := json.Unmarshal(arg, &nodes); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("err: %s", err))
		os.Exit(1)
	}

	// Build a map of unique node specifications (cpu, ram, gpu) and assign of weight of 1
	// to each specifications without weight or with a weight smaller than 1.
	SpecsWeightsMap := make(map[NodeSpecs]int)
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Specs.Weight < 1 {
			SpecsWeightsMap[nodes[i].Specs] = 1
		}
	}

	// Extract list of unique specifications from map
	UniqueSpecs := make([]NodeSpecs, len(SpecsWeightsMap))
	i := 0
	for key := range SpecsWeightsMap {
		UniqueSpecs[i] = key
		i++
	}

	// Sort the specifications based on the importance of each components GPUs > RealMemory > CPUs
	sort.SliceStable(UniqueSpecs, func(i, j int) bool {
		if UniqueSpecs[i].Gpus != UniqueSpecs[j].Gpus {
			return UniqueSpecs[i].Gpus < UniqueSpecs[j].Gpus
		} else if UniqueSpecs[i].RealMemory != UniqueSpecs[j].RealMemory {
			return UniqueSpecs[i].RealMemory < UniqueSpecs[j].RealMemory
		}
		return UniqueSpecs[i].Cpus < UniqueSpecs[j].Cpus
	})

	// Assign a weight to each specifications based on its rank
	for i := 0; i < len(UniqueSpecs); i++ {
		SpecsWeightsMap[UniqueSpecs[i]] = i + 1
	}

	// Assign the weight to the nodes
	// Also identify node prefix and index based on node name
	for i := 0; i < len(nodes); i++ {
		if nodes[i].Specs.Weight < 1 {
			nodes[i].Specs.Weight = SpecsWeightsMap[nodes[i].Specs]
		}
		nodes[i].Prefix = strings.TrimRightFunc(nodes[i].Name, unicode.IsNumber)
		nodes[i].Index, _ = strconv.Atoi(strings.TrimPrefix(nodes[i].Name, nodes[i].Prefix))
	}

	// Sort the nodes based on the weight of their specifications
	sort.SliceStable(nodes, func(i, j int) bool {
		if nodes[i].Specs.Weight != nodes[j].Specs.Weight {
			return nodes[i].Specs.Weight < nodes[j].Specs.Weight
		} else if nodes[i].Prefix != nodes[j].Prefix {
			return nodes[i].Prefix < nodes[j].Prefix
		} else {
			return nodes[i].Index < nodes[j].Index
		}
	})

	// Output the nodes in Slurm configuration format with the weights
	for i := 0; i < len(nodes); i++ {
		fmt.Printf(
			"NodeName=%s CPUs=%d RealMemory=%d MemSpecLimit=%d ",
			nodes[i].Name,
			nodes[i].Specs.Cpus,
			nodes[i].Specs.RealMemory,
			nodes[i].Specs.MemSpecLimit,
		)
		if nodes[i].Specs.Gpus > 0 {
			fmt.Printf("Gres=gpu:%d ", nodes[i].Specs.Gpus)
		}
		fmt.Printf("Weight=%d\n", nodes[i].Specs.Weight)
	}
	os.Exit(0)
}
