package tags

import (
	"encoding/binary"
	"encoding/hex"
	"strings"
)

// ExtractUPIDFromSCTE35Hex decodes a hexadecimal SCTE-35 splice_info_section
// and extracts the raw segmentation_upid() data from the first segmentation descriptor found.
// Returns the UPID value as a string (which may contain any bytes).
func ExtractUPIDFromSCTE35Hex(scte35Hex string) (string, bool) {
	payload := strings.TrimSpace(scte35Hex)
	payload = strings.TrimPrefix(payload, "0x")
	payload = strings.TrimPrefix(payload, "0X")

	data, err := hex.DecodeString(payload)
	if err != nil || len(data) < 14 {
		return "", false
	}

	// Extract splice_command_length (lower 12 bits at bytes 11-12)
	cmdLen := int(binary.BigEndian.Uint16(data[11:13]) & 0x0FFF)
	if cmdLen < 0 || 14+cmdLen > len(data) {
		return "", false
	}

	// Extract descriptor_loop_length (16 bits after splice_command)
	return extractUPIDFromDescriptorLoop(data, 14+cmdLen)
}

// extractUPIDFromDescriptorLoop navigates the descriptor loop and extracts the first UPID.
func extractUPIDFromDescriptorLoop(data []byte, loopStart int) (string, bool) {
	if loopStart+2 > len(data) {
		return "", false
	}

	loopLen := int(binary.BigEndian.Uint16(data[loopStart : loopStart+2]))
	loopEnd := loopStart + 2 + loopLen
	if loopEnd > len(data) {
		return "", false
	}

	// Iterate through splice_descriptors: each has tag(1) + length(1) + data
	pos := loopStart + 2
	for pos+2 <= loopEnd && pos+2 <= len(data) {
		tag := data[pos]
		descLen := int(data[pos+1])
		descStart := pos + 2
		descEnd := descStart + descLen

		if descEnd > len(data) || descEnd > loopEnd {
			return "", false
		}

		// 0x02 is the SegmentationDescriptor tag
		if tag == 0x02 {
			if upid, ok := parseSegmentationDescriptor(data[descStart:descEnd]); ok {
				return upid, true
			}
		}

		pos = descEnd
	}

	return "", false
}

// parseSegmentationDescriptor extracts the segmentation_upid() from a descriptor.
// The descriptor structure (from SCTE-35 spec):
//   - identifier (4 bytes, "CUEI")
//   - segmentation_event_id (4 bytes)
//   - segmentation_event_cancel_indicator + flags (1 byte)
//   - additional flag fields depending on flags
//   - segmentation_upid: type(1) + length(1) + data
func parseSegmentationDescriptor(d []byte) (string, bool) {
	// Minimum size: identifier(4) + event_id(4) + flags(1) + upid_type(1) + upid_length(1)
	if len(d) < 11 {
		return "", false
	}

	// Skip identifier (4 bytes) and segmentation_event_id (4 bytes)
	pos := 8

	// Check segmentation_event_cancel_indicator (bit 7)
	if d[pos]&0x80 != 0 {
		// Event is cancelled; no UPID data
		return "", false
	}
	pos++

	// Parse remaining flags and skip optional fields
	pos = skipOptionalSegmentationFields(d, pos)
	if pos < 0 {
		return "", false
	}

	// Read segmentation_upid: type(1) + length(1) + data
	if pos+2 > len(d) {
		return "", false
	}

	upidLen := int(d[pos+1])
	upidStart := pos + 2
	upidEnd := upidStart + upidLen

	if upidLen == 0 || upidEnd > len(d) {
		return "", false
	}

	return string(d[upidStart:upidEnd]), true
}

// skipOptionalSegmentationFields handles variable-length fields in the segmentation descriptor.
// Returns the position after all optional fields, or -1 on error.
func skipOptionalSegmentationFields(d []byte, pos int) int {
	if pos >= len(d) {
		return -1
	}

	flags := d[pos]

	// program_segmentation_flag (bit 7): if 0, components follow
	if flags&0x80 == 0 {
		if pos+1 >= len(d) {
			return -1
		}
		componentCount := int(d[pos+1])
		pos += 2 + (componentCount * 6) // Each component is 6 bytes
	} else {
		pos++ // Skip the flags byte for program splice
	}

	// segmentation_duration_flag (bit 6): if 1, 5 bytes follow
	if flags&0x40 != 0 {
		if pos+5 > len(d) {
			return -1
		}
		pos += 5
	}

	return pos
}
