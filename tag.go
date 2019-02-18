package govalid

import (
	"strconv"
	"strings"
)

func getBoolFromTag(tag string, key string) bool {
	cs := strings.Split(tag, ",")
	for _, c := range cs {
		if c == key {
			return true
		}
	}
	return false
}

func getIntFromTag(tag string, key string) int {
	cs := strings.Split(tag, ",")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			i, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	return 0
}

func getInt64FromTag(tag string, key string) int64 {
	cs := strings.Split(tag, ",")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			i, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				panic(err)
			}
			return i
		}
	}
	return 0
}

func getFloat32FromTag(tag string, key string) float32 {
	cs := strings.Split(tag, ",")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			f, err := strconv.ParseFloat(parts[1], 32)
			if err != nil {
				panic(err)
			}
			return float32(f)
		}
	}
	return 0
}

func getFloat64FromTag(tag string, key string) float64 {
	cs := strings.Split(tag, ",")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			f, err := strconv.ParseFloat(parts[1], 64)
			if err != nil {
				panic(err)
			}
			return f
		}
	}
	return 0
}

func getStringFromTag(tag string, key string) string {
	cs := strings.Split(tag, ",")
	for _, c := range cs {
		parts := strings.Split(c, ":")
		if len(parts) != 2 {
			continue
		}
		if parts[0] == key {
			return parts[1]
		}
	}
	return ""
}
