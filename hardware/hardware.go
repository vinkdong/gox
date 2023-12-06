package hardware

const (
	MaxCpuCount   = 64
	MaxMemorySize = 512
)

// calculateMaxMemory 根据CPU核心数计算最大内存。
func calculateMaxMemory(cpuCores uint32) uint32 {
	// 假设每个CPU核心最多对应8GB内存
	if cpuCores >= 2 {
		return cpuCores * 8
	}
	return cpuCores * 2
}

// calculateRequiredMemory 根据CPU核心数和所需内存计算出符合条件的最小内存。
func calculateRequiredMemory(cpuCores uint32, memoryGB uint32) uint32 {
	for _, v := range []uint32{2, 4, 8} {
		if m := cpuCores * v; m >= memoryGB {
			return m
		}
	}
	return 0
}

// findValidCPUCore 查找符合规则的下一个有效CPU核心数。
func findValidCPUCore(cpuCores uint32) uint32 {
	// 特定的CPU核心数列表
	validCores := []uint32{1, 2, 4, 8, 12, 16, 24, 32, 64}

	for _, core := range validCores {
		if core >= cpuCores {
			return core
		}
	}

	// 如果超出范围，返回最大核心数
	return 64
}

// StandardizeConfig 根据CPU核心数和内存需求，给出一个合理的配置
func StandardizeConfig(cpuCores uint32, memoryGB uint32, maxCpuCountV ...uint32) (uint32, uint32, uint32) {
	maxCpuCount := uint32(MaxCpuCount)
	maxMemorySize := uint32(MaxMemorySize)
	if len(maxCpuCountV) > 0 && maxCpuCountV[0] > 0 {
		maxCpuCount = maxCpuCountV[0]
		maxMemorySize = calculateMaxMemory(maxCpuCount)
	}
	cpuSplitQuantity := IntDivCeil(cpuCores, maxCpuCount)
	memorySplitQuantity := IntDivCeil(memoryGB, maxMemorySize)
	quantity := cpuSplitQuantity
	if cpuSplitQuantity < memorySplitQuantity {
		quantity = memorySplitQuantity
	}

	cpuCores = IntDivCeil(cpuCores, quantity)
	memoryGB = IntDivCeil(memoryGB, quantity)

	cpuCores = findValidCPUCore(cpuCores)
	maxMemory := calculateMaxMemory(cpuCores)

	// 增加CPU核心数，直到找到满足内存需求的配置
	for maxMemory < memoryGB {
		cpuCores = findValidCPUCore(cpuCores + 1)
		maxMemory = calculateMaxMemory(cpuCores)
	}

	return cpuCores, calculateRequiredMemory(cpuCores, memoryGB), quantity
}

func IntDivCeil(a, b uint32) uint32 {
	c := a / b
	if a%b != 0 {
		c++
	}
	return c
}
