/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package state

import (
	"encoding/json"

	"github.com/colo-sh/cpumanager/checkpointmanager"
)


type Checksum uint64
var _ checkpointmanager.Checkpoint = &CPUManagerCheckpointV1{}
var _ checkpointmanager.Checkpoint = &CPUManagerCheckpointV2{}
var _ checkpointmanager.Checkpoint = &CPUManagerCheckpoint{}

// CPUManagerCheckpoint struct is used to store cpu/pod assignments in a checkpoint in v2 format
type CPUManagerCheckpoint struct {
	PolicyName    string                       `json:"policyName"`
	DefaultCPUSet string                       `json:"defaultCpuSet"`
	Entries       map[string]map[string]string `json:"entries,omitempty"`
	Checksum      Checksum            `json:"checksum"`
}

// CPUManagerCheckpointV1 struct is used to store cpu/pod assignments in a checkpoint in v1 format
type CPUManagerCheckpointV1 struct {
	PolicyName    string            `json:"policyName"`
	DefaultCPUSet string            `json:"defaultCpuSet"`
	Entries       map[string]string `json:"entries,omitempty"`
	Checksum      Checksum `json:"checksum"`
}

// CPUManagerCheckpointV2 struct is used to store cpu/pod assignments in a checkpoint in v2 format
type CPUManagerCheckpointV2 = CPUManagerCheckpoint

// NewCPUManagerCheckpoint returns an instance of Checkpoint
func NewCPUManagerCheckpoint() *CPUManagerCheckpoint {
	//nolint:staticcheck // unexported-type-in-api user-facing error message
	return newCPUManagerCheckpointV2()
}

func newCPUManagerCheckpointV1() *CPUManagerCheckpointV1 {
	return &CPUManagerCheckpointV1{
		Entries: make(map[string]string),
	}
}

func newCPUManagerCheckpointV2() *CPUManagerCheckpointV2 {
	return &CPUManagerCheckpointV2{
		Entries: make(map[string]map[string]string),
	}
}

// MarshalCheckpoint returns marshalled checkpoint in v1 format
func (cp *CPUManagerCheckpointV1) MarshalCheckpoint() ([]byte, error) {
	// make sure checksum wasn't set before so it doesn't affect output checksum
	//cp.Checksum = 0
	//cp.Checksum = checksum.New(cp)
	return json.Marshal(*cp)
}

// MarshalCheckpoint returns marshalled checkpoint in v2 format
func (cp *CPUManagerCheckpointV2) MarshalCheckpoint() ([]byte, error) {
	// make sure checksum wasn't set before so it doesn't affect output checksum
	//cp.Checksum = 0
	//cp.Checksum = checksum.New(cp)
	return json.Marshal(*cp)
}

// UnmarshalCheckpoint tries to unmarshal passed bytes to checkpoint in v1 format
func (cp *CPUManagerCheckpointV1) UnmarshalCheckpoint(blob []byte) error {
	return json.Unmarshal(blob, cp)
}

// UnmarshalCheckpoint tries to unmarshal passed bytes to checkpoint in v2 format
func (cp *CPUManagerCheckpointV2) UnmarshalCheckpoint(blob []byte) error {
	return json.Unmarshal(blob, cp)
}
