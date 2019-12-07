package ai

import (
	"strconv"
	"strings"
)

// Descriptor ...
type Descriptor struct {
	Name        string
	VersionFrom int
	VersionTo   int
}

// DescriptorFromString constructs a descriptor from its string representation
func DescriptorFromString(s string) Descriptor {
	descriptorSplit := strings.Split(s, ":")

	descriptor := Descriptor{}

	descriptor.Name = descriptorSplit[0]

	if len(descriptorSplit) == 1 {
		descriptor.VersionFrom = -1
		descriptor.VersionTo = -1
	} else {
		rangeDescriptorSplit := strings.Split(descriptorSplit[1], "..")

		switch len(rangeDescriptorSplit) {
		case 1:
			if len(rangeDescriptorSplit[0]) > 0 {
				descriptorVersion, _ := strconv.Atoi(rangeDescriptorSplit[0])

				descriptor.VersionFrom = descriptorVersion
				descriptor.VersionTo = descriptorVersion
			} else {
				descriptor.VersionFrom = 0
				descriptor.VersionTo = -1
			}

		case 2:
			descriptor.VersionFrom = 0
			if len(rangeDescriptorSplit[0]) > 0 {
				from, _ := strconv.Atoi(rangeDescriptorSplit[0])
				descriptor.VersionFrom = from
			}

			descriptor.VersionTo = -1
			if len(rangeDescriptorSplit[1]) > 0 {
				to, _ := strconv.Atoi(rangeDescriptorSplit[1])
				descriptor.VersionTo = to
			}
		}
	}

	return descriptor
}
