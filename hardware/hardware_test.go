package hardware

import "testing"

func TestStandardizeConfig(t *testing.T) {
	testCases := []struct {
		cpuCores     uint32
		memoryGB     uint32
		maxCpuCores  uint32
		wantCpuCores uint32
		wantMemoryGB uint32
		wantQuality  uint32
	}{
		{12, 513, 8, 8, 64, 9},
		{4, 32, 0, 4, 32, 1},
		{2, 64, 0, 8, 64, 1},
		{8, 65, 4, 4, 32, 3},
	}

	for _, tc := range testCases {
		gotCpuCores, gotMemoryGB, gotQuality := StandardizeConfig(tc.cpuCores, tc.memoryGB, tc.maxCpuCores)
		if gotCpuCores != tc.wantCpuCores || gotMemoryGB != tc.wantMemoryGB || gotQuality != tc.wantQuality {
			t.Errorf("StandardizeConfig(%d, %d, %d) = %d, %d, %d; want %d, %d, %d",
				tc.cpuCores, tc.memoryGB, tc.maxCpuCores, gotCpuCores, gotMemoryGB, gotQuality, tc.wantCpuCores, tc.wantMemoryGB, tc.wantQuality)
		}
	}
}
