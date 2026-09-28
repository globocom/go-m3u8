package playlist_test

import (
	"os"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
	"github.com/stretchr/testify/assert"
)

// Master Manifest filters
func TestVariantsHeightFilter(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMaxHeight(720)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 7)
}

func TestVariantsMinHeightFilter(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMinHeight(720)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 3)
}

func TestVariantsMaxWidthFilter(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMaxWidth(1280)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 7)
}

func TestVariantsMinWidthFilter(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMinWidth(768)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 4)
}

func TestVariantsMixedFiltersMinHeightAndMaxWidth(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMinHeight(720)
	playlist.FilterByMaxWidth(1280)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 2)
}

func TestVariantsMixedFiltersMinWidthMaxWidthAndMinHeight(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMinWidth(500)
	playlist.FilterByMaxWidth(1000)
	playlist.FilterByMinHeight(300)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 2)
}

func TestVariantsMaxBitrateFilter(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMaxBitrate(1000000)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 4)
}

func TestVariantsMinBitrateFilter(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMinBitrate(1000000)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 4)
}

func TestVariantsMixedFiltersMinAndMaxBitrate(t *testing.T) {
	file, _ := os.Open("./../mocks/multivariant/multivariant.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByMinBitrate(700000)
	playlist.FilterByMaxBitrate(2000000)
	nodes := playlist.Variants()
	assert.NotNil(t, nodes)
	assert.Len(t, nodes, 3)
}

// Media Manifest filters
func TestFilterByDVRWindow(t *testing.T) {
	file, _ := os.Open("./../mocks/media/media.m3u8")
	playlist, err := m3u8.ParsePlaylist(file)
	assert.NoError(t, err)

	playlist.FilterByDVRWindow(20)
	segments := playlist.Segments()
	assert.NotNil(t, segments)
	assert.Len(t, segments, 4)
	assert.Equal(t, "channel-audio_1=96000-video=3442944-364042192.ts", segments[0].HLSElement.URI)
	assert.Equal(t, "channel-audio_1=96000-video=3442944-364042195.ts", segments[3].HLSElement.URI)
	assert.Equal(t, "364042192", playlist.MediaSequenceValue())
	assert.Equal(t, 364042192, playlist.MediaSequence)
	assert.Equal(t, 19.2, playlist.DVR)
}
