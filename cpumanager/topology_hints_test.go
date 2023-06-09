/*
Copyright 2019 The Kubernetes Authors.

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

package cpumanager

import (
	"reflect"
	"sort"
	"testing"

	cadvisorapi "github.com/google/cadvisor/info/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"github.com/colo-sh/cpumanager/cpumanager/state"
	"github.com/colo-sh/cpumanager/cpumanager/topology"
	"github.com/colo-sh/cpumanager/cpuset"
	"github.com/colo-sh/cpumanager/topologymanager"
	"github.com/colo-sh/cpumanager/topologymanager/bitmask"
)

type testCase struct {
	name          string
	pod           v1.Pod
	container     v1.Container
	assignments   state.ContainerCPUAssignments
	defaultCPUSet cpuset.CPUSet
	expectedHints []topologymanager.TopologyHint
}

func returnMachineInfo() cadvisorapi.MachineInfo {
	return cadvisorapi.MachineInfo{
		NumCores: 12,
		Topology: []cadvisorapi.Node{
			{Id: 0,
				Cores: []cadvisorapi.Core{
					{SocketID: 0, Id: 0, Threads: []int{0, 6}},
					{SocketID: 0, Id: 1, Threads: []int{1, 7}},
					{SocketID: 0, Id: 2, Threads: []int{2, 8}},
				},
			},
			{Id: 1,
				Cores: []cadvisorapi.Core{
					{SocketID: 1, Id: 0, Threads: []int{3, 9}},
					{SocketID: 1, Id: 1, Threads: []int{4, 10}},
					{SocketID: 1, Id: 2, Threads: []int{5, 11}},
				},
			},
		},
	}
}

func TestPodGuaranteedCPUs(t *testing.T) {
	CPUs := [][]struct {
		request string
		limit   string
	}{
		{
			{request: "0", limit: "0"},
		},
		{
			{request: "2", limit: "2"},
		},
		{
			{request: "5", limit: "5"},
		},
		{
			{request: "2", limit: "2"},
			{request: "4", limit: "4"},
		},
	}
	// tc for not guaranteed Pod
	testPod1 := makeMultiContainerPod(CPUs[0], CPUs[0])
	testPod2 := makeMultiContainerPod(CPUs[0], CPUs[1])
	testPod3 := makeMultiContainerPod(CPUs[1], CPUs[0])
	// tc for guaranteed Pod
	testPod4 := makeMultiContainerPod(CPUs[1], CPUs[1])
	testPod5 := makeMultiContainerPod(CPUs[2], CPUs[2])
	// tc for comparing init containers and user containers
	testPod6 := makeMultiContainerPod(CPUs[1], CPUs[2])
	testPod7 := makeMultiContainerPod(CPUs[2], CPUs[1])
	// tc for multi containers
	testPod8 := makeMultiContainerPod(CPUs[3], CPUs[3])

	p := staticPolicy{}

	tcases := []struct {
		name        string
		pod         *v1.Pod
		expectedCPU int
	}{
		{
			name:        "TestCase01: if requestedCPU == 0, Pod is not Guaranteed Qos",
			pod:         testPod1,
			expectedCPU: 0,
		},
		{
			name:        "TestCase02: if requestedCPU == 0, Pod is not Guaranteed Qos",
			pod:         testPod2,
			expectedCPU: 0,
		},
		{
			name:        "TestCase03: if requestedCPU == 0, Pod is not Guaranteed Qos",
			pod:         testPod3,
			expectedCPU: 0,
		},
		{
			name:        "TestCase04: Guaranteed Pod requests 2 CPUs",
			pod:         testPod4,
			expectedCPU: 2,
		},
		{
			name:        "TestCase05: Guaranteed Pod requests 5 CPUs",
			pod:         testPod5,
			expectedCPU: 5,
		},
		{
			name:        "TestCase06: The number of CPUs requested By app is bigger than the number of CPUs requested by init",
			pod:         testPod6,
			expectedCPU: 5,
		},
		{
			name:        "TestCase07: The number of CPUs requested By init is bigger than the number of CPUs requested by app",
			pod:         testPod7,
			expectedCPU: 5,
		},
		{
			name:        "TestCase08: Sum of CPUs requested by multiple containers",
			pod:         testPod8,
			expectedCPU: 6,
		},
	}
	for _, tc := range tcases {
		requestedCPU := p.podGuaranteedCPUs(tc.pod)

		if requestedCPU != tc.expectedCPU {
			t.Errorf("Expected in result to be %v , got %v", tc.expectedCPU, requestedCPU)
		}
	}
}

func TestGetTopologyHints(t *testing.T) {
	machineInfo := returnMachineInfo()
	tcases := returnTestCases()

	for _, tc := range tcases {
		topology, _ := topology.Discover(&machineInfo)

		var activePods []*v1.Pod
		for p := range tc.assignments {
			pod := v1.Pod{}
			pod.UID = types.UID(p)
			for c := range tc.assignments[p] {
				container := v1.Container{}
				container.Name = c
				pod.Spec.Containers = append(pod.Spec.Containers, container)
			}
			activePods = append(activePods, &pod)
		}

		m := manager{
			policy: &staticPolicy{
				topology: topology,
			},
			state: &mockState{
				assignments:   tc.assignments,
				defaultCPUSet: tc.defaultCPUSet,
			},
			topology:          topology,
			activePods:        func() []*v1.Pod { return activePods },
			podStatusProvider: mockPodStatusProvider{},
			sourcesReady:      &sourcesReadyStub{},
		}

		hints := m.GetTopologyHints(&tc.pod, &tc.container)[string(v1.ResourceCPU)]
		if len(tc.expectedHints) == 0 && len(hints) == 0 {
			continue
		}

		if m.pendingAdmissionPod == nil {
			t.Errorf("The pendingAdmissionPod should point to the current pod after the call to GetTopologyHints()")
		}

		sort.SliceStable(hints, func(i, j int) bool {
			return hints[i].LessThan(hints[j])
		})
		sort.SliceStable(tc.expectedHints, func(i, j int) bool {
			return tc.expectedHints[i].LessThan(tc.expectedHints[j])
		})
		if !reflect.DeepEqual(tc.expectedHints, hints) {
			t.Errorf("Expected in result to be %v , got %v", tc.expectedHints, hints)
		}
	}
}

func TestGetPodTopologyHints(t *testing.T) {
	machineInfo := returnMachineInfo()

	for _, tc := range returnTestCases() {
		topology, _ := topology.Discover(&machineInfo)

		var activePods []*v1.Pod
		for p := range tc.assignments {
			pod := v1.Pod{}
			pod.UID = types.UID(p)
			for c := range tc.assignments[p] {
				container := v1.Container{}
				container.Name = c
				pod.Spec.Containers = append(pod.Spec.Containers, container)
			}
			activePods = append(activePods, &pod)
		}

		m := manager{
			policy: &staticPolicy{
				topology: topology,
			},
			state: &mockState{
				assignments:   tc.assignments,
				defaultCPUSet: tc.defaultCPUSet,
			},
			topology:          topology,
			activePods:        func() []*v1.Pod { return activePods },
			podStatusProvider: mockPodStatusProvider{},
			sourcesReady:      &sourcesReadyStub{},
		}

		podHints := m.GetPodTopologyHints(&tc.pod)[string(v1.ResourceCPU)]
		if len(tc.expectedHints) == 0 && len(podHints) == 0 {
			continue
		}

		sort.SliceStable(podHints, func(i, j int) bool {
			return podHints[i].LessThan(podHints[j])
		})
		sort.SliceStable(tc.expectedHints, func(i, j int) bool {
			return tc.expectedHints[i].LessThan(tc.expectedHints[j])
		})
		if !reflect.DeepEqual(tc.expectedHints, podHints) {
			t.Errorf("Expected in result to be %v , got %v", tc.expectedHints, podHints)
		}
	}
}

func returnTestCases() []testCase {
	testPod1 := makePod("fakePod", "fakeContainer", "2", "2")
	testContainer1 := &testPod1.Spec.Containers[0]
	testPod2 := makePod("fakePod", "fakeContainer", "5", "5")
	testContainer2 := &testPod2.Spec.Containers[0]
	testPod3 := makePod("fakePod", "fakeContainer", "7", "7")
	testContainer3 := &testPod3.Spec.Containers[0]
	testPod4 := makePod("fakePod", "fakeContainer", "11", "11")
	testContainer4 := &testPod4.Spec.Containers[0]

	firstSocketMask, _ := bitmask.NewBitMask(0)
	secondSocketMask, _ := bitmask.NewBitMask(1)
	crossSocketMask, _ := bitmask.NewBitMask(0, 1)

	return []testCase{
		{
			name:          "Request 2 CPUs, 4 available on NUMA 0, 6 available on NUMA 1",
			pod:           *testPod1,
			container:     *testContainer1,
			defaultCPUSet: cpuset.NewCPUSet(2, 3, 4, 5, 6, 7, 8, 9, 10, 11),
			expectedHints: []topologymanager.TopologyHint{
				{
					NUMANodeAffinity: firstSocketMask,
					Preferred:        true,
				},
				{
					NUMANodeAffinity: secondSocketMask,
					Preferred:        true,
				},
				{
					NUMANodeAffinity: crossSocketMask,
					Preferred:        false,
				},
			},
		},
		{
			name:          "Request 5 CPUs, 4 available on NUMA 0, 6 available on NUMA 1",
			pod:           *testPod2,
			container:     *testContainer2,
			defaultCPUSet: cpuset.NewCPUSet(2, 3, 4, 5, 6, 7, 8, 9, 10, 11),
			expectedHints: []topologymanager.TopologyHint{
				{
					NUMANodeAffinity: secondSocketMask,
					Preferred:        true,
				},
				{
					NUMANodeAffinity: crossSocketMask,
					Preferred:        false,
				},
			},
		},
		{
			name:          "Request 7 CPUs, 4 available on NUMA 0, 6 available on NUMA 1",
			pod:           *testPod3,
			container:     *testContainer3,
			defaultCPUSet: cpuset.NewCPUSet(2, 3, 4, 5, 6, 7, 8, 9, 10, 11),
			expectedHints: []topologymanager.TopologyHint{
				{
					NUMANodeAffinity: crossSocketMask,
					Preferred:        true,
				},
			},
		},
		{
			name:          "Request 11 CPUs, 4 available on NUMA 0, 6 available on NUMA 1",
			pod:           *testPod4,
			container:     *testContainer4,
			defaultCPUSet: cpuset.NewCPUSet(2, 3, 4, 5, 6, 7, 8, 9, 10, 11),
			expectedHints: nil,
		},
		{
			name:          "Request 2 CPUs, 1 available on NUMA 0, 1 available on NUMA 1",
			pod:           *testPod1,
			container:     *testContainer1,
			defaultCPUSet: cpuset.NewCPUSet(0, 3),
			expectedHints: []topologymanager.TopologyHint{
				{
					NUMANodeAffinity: crossSocketMask,
					Preferred:        false,
				},
			},
		},
		{
			name:          "Request more CPUs than available",
			pod:           *testPod2,
			container:     *testContainer2,
			defaultCPUSet: cpuset.NewCPUSet(0, 1, 2, 3),
			expectedHints: nil,
		},
		{
			name:      "Regenerate Single-Node NUMA Hints if already allocated 1/2",
			pod:       *testPod1,
			container: *testContainer1,
			assignments: state.ContainerCPUAssignments{
				string(testPod1.UID): map[string]cpuset.CPUSet{
					testContainer1.Name: cpuset.NewCPUSet(0, 6),
				},
			},
			defaultCPUSet: cpuset.NewCPUSet(),
			expectedHints: []topologymanager.TopologyHint{
				{
					NUMANodeAffinity: firstSocketMask,
					Preferred:        true,
				},
				{
					NUMANodeAffinity: crossSocketMask,
					Preferred:        false,
				},
			},
		},
		{
			name:      "Regenerate Single-Node NUMA Hints if already allocated 1/2",
			pod:       *testPod1,
			container: *testContainer1,
			assignments: state.ContainerCPUAssignments{
				string(testPod1.UID): map[string]cpuset.CPUSet{
					testContainer1.Name: cpuset.NewCPUSet(3, 9),
				},
			},
			defaultCPUSet: cpuset.NewCPUSet(),
			expectedHints: []topologymanager.TopologyHint{
				{
					NUMANodeAffinity: secondSocketMask,
					Preferred:        true,
				},
				{
					NUMANodeAffinity: crossSocketMask,
					Preferred:        false,
				},
			},
		},
		{
			name:      "Regenerate Cross-NUMA Hints if already allocated",
			pod:       *testPod4,
			container: *testContainer4,
			assignments: state.ContainerCPUAssignments{
				string(testPod4.UID): map[string]cpuset.CPUSet{
					testContainer4.Name: cpuset.NewCPUSet(0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10),
				},
			},
			defaultCPUSet: cpuset.NewCPUSet(),
			expectedHints: []topologymanager.TopologyHint{
				{
					NUMANodeAffinity: crossSocketMask,
					Preferred:        true,
				},
			},
		},
		{
			name:      "Requested less than already allocated",
			pod:       *testPod1,
			container: *testContainer1,
			assignments: state.ContainerCPUAssignments{
				string(testPod1.UID): map[string]cpuset.CPUSet{
					testContainer1.Name: cpuset.NewCPUSet(0, 6, 3, 9),
				},
			},
			defaultCPUSet: cpuset.NewCPUSet(),
			expectedHints: []topologymanager.TopologyHint{},
		},
		{
			name:      "Requested more than already allocated",
			pod:       *testPod4,
			container: *testContainer4,
			assignments: state.ContainerCPUAssignments{
				string(testPod4.UID): map[string]cpuset.CPUSet{
					testContainer4.Name: cpuset.NewCPUSet(0, 6, 3, 9),
				},
			},
			defaultCPUSet: cpuset.NewCPUSet(),
			expectedHints: []topologymanager.TopologyHint{},
		},
	}
}
