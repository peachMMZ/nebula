package util

import (
	"strconv"
	"strings"
)

// CompareVersion 比较两个语义化版本号
// 返回值：1 表示 v1 > v2，-1 表示 v1 < v2，0 表示 v1 == v2
// 支持格式：1.0.0, v1.0.0, 1.0.0-beta.1, 1.0.0+build.123
func CompareVersion(v1, v2 string) int {
	// 移除 v 前缀
	v1 = strings.TrimPrefix(v1, "v")
	v2 = strings.TrimPrefix(v2, "v")

	// 分离主版本号和预发布/构建元数据
	v1Main, v1Pre := splitVersion(v1)
	v2Main, v2Pre := splitVersion(v2)

	// 比较主版本号
	result := compareMainVersion(v1Main, v2Main)
	if result != 0 {
		return result
	}

	// 主版本号相同，比较预发布版本
	return comparePreRelease(v1Pre, v2Pre)
}

// splitVersion 分离主版本号和预发布标识
func splitVersion(version string) (main, preRelease string) {
	// 移除构建元数据 (+build.xxx)
	if idx := strings.Index(version, "+"); idx != -1 {
		version = version[:idx]
	}

	// 分离预发布标识 (-beta.1)
	if idx := strings.Index(version, "-"); idx != -1 {
		return version[:idx], version[idx+1:]
	}

	return version, ""
}

// compareMainVersion 比较主版本号 (major.minor.patch)
func compareMainVersion(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	// 补齐到相同长度
	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var num1, num2 int

		if i < len(parts1) {
			num1, _ = strconv.Atoi(parts1[i])
		}
		if i < len(parts2) {
			num2, _ = strconv.Atoi(parts2[i])
		}

		if num1 > num2 {
			return 1
		}
		if num1 < num2 {
			return -1
		}
	}

	return 0
}

// comparePreRelease 比较预发布版本
// 没有预发布标识的版本 > 有预发布标识的版本
// 例如：1.0.0 > 1.0.0-beta
func comparePreRelease(pre1, pre2 string) int {
	// 都没有预发布标识，相等
	if pre1 == "" && pre2 == "" {
		return 0
	}

	// v1 是正式版，v2 是预发布版
	if pre1 == "" && pre2 != "" {
		return 1
	}

	// v1 是预发布版，v2 是正式版
	if pre1 != "" && pre2 == "" {
		return -1
	}

	// 都是预发布版，按字典序比较
	if pre1 < pre2 {
		return -1
	}
	if pre1 > pre2 {
		return 1
	}

	return 0
}

// IsNewerVersion 检查 newVer 是否比 currentVer 新
func IsNewerVersion(currentVer, newVer string) bool {
	return CompareVersion(newVer, currentVer) > 0
}
