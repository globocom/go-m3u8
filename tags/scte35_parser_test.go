package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractUPIDFromSCTE35Hex_ShortValue(t *testing.T) {
	// SCTE-35 with segmentation_upid: type=0x08 (TI), length=0x04, data="2 PT" (0x32205054)
	scte35Hex := "0xFC303F00000000000000FFF01405C00000007FEFFFDB6680287E00A55820B7190002001A0218435545490000016B7FFF0000A5584408043220505400000047F03520"
	upid, ok := ExtractUPIDFromSCTE35Hex(scte35Hex)
	assert.True(t, ok)
	assert.Equal(t, "2 PT", upid)
}

func TestExtractUPIDFromSCTE35Hex_LongValue(t *testing.T) {
	// SCTE-35 with segmentation_upid: type=0x0C (MPU), length=0x0F (15 bytes), data=" RED2026091585 "
	scte35Hex := "0xFC304A000000000BB800FFF01405C00000007FEFFE3E2E4620FE0036EE804414000000250223435545490000049B7FFF000036EE800C0F205245443230323630393135383520100101428611E9"
	upid, ok := ExtractUPIDFromSCTE35Hex(scte35Hex)
	assert.True(t, ok)
	assert.Equal(t, " RED2026091585 ", upid)
}

func TestExtractUPIDFromSCTE35Hex_NoSegmentationDescriptor(t *testing.T) {
	// SCTE-35 without segmentation descriptor (descriptor_loop_length = 0)
	scte35Hex := "0xFC3025000000000BB800FFF01405F00001BB7FEFFE06EF5210FE005265C0000101010000E50D79A2"
	upid, ok := ExtractUPIDFromSCTE35Hex(scte35Hex)
	assert.False(t, ok)
	assert.Equal(t, "", upid)
}

func TestExtractUPIDFromSCTE35Hex_InvalidHexadecimal(t *testing.T) {
	// Invalid hex string
	scte35Hex := "0xZZZZ"
	upid, ok := ExtractUPIDFromSCTE35Hex(scte35Hex)
	assert.False(t, ok)
	assert.Equal(t, "", upid)
}

func TestExtractUPIDFromSCTE35Hex_TooShortData(t *testing.T) {
	// Too short to be a valid SCTE-35 section
	scte35Hex := "0xFC30"
	upid, ok := ExtractUPIDFromSCTE35Hex(scte35Hex)
	assert.False(t, ok)
	assert.Equal(t, "", upid)
}

func TestExtractUPIDFromSCTE35Hex_WithPrefixVariations(t *testing.T) {
	// Test with lowercase 0x prefix
	upid1, ok1 := ExtractUPIDFromSCTE35Hex("0xFC303F00000000000000FFF01405C00000007FEFFFDB6680287E00A55820B7190002001A0218435545490000016B7FFF0000A5584408043220505400000047F03520")
	assert.True(t, ok1)
	assert.Equal(t, "2 PT", upid1)

	// Test with uppercase 0X prefix
	upid2, ok2 := ExtractUPIDFromSCTE35Hex("0XFC303F00000000000000FFF01405C00000007FEFFFDB6680287E00A55820B7190002001A0218435545490000016B7FFF0000A5584408043220505400000047F03520")
	assert.True(t, ok2)
	assert.Equal(t, "2 PT", upid2)

	// Test without prefix
	upid3, ok3 := ExtractUPIDFromSCTE35Hex("FC303F00000000000000FFF01405C00000007FEFFFDB6680287E00A55820B7190002001A0218435545490000016B7FFF0000A5584408043220505400000047F03520")
	assert.True(t, ok3)
	assert.Equal(t, "2 PT", upid3)
}

func TestExtractUPIDFromSCTE35Hex_WithWhitespace(t *testing.T) {
	// Test with leading/trailing whitespace
	upid, ok := ExtractUPIDFromSCTE35Hex("  0xFC303F00000000000000FFF01405C00000007FEFFFDB6680287E00A55820B7190002001A0218435545490000016B7FFF0000A5584408043220505400000047F03520  ")
	assert.True(t, ok)
	assert.Equal(t, "2 PT", upid)
}
