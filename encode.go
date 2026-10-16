package go_m3u8

import (
	"fmt"
	"strings"

	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/globocom/go-m3u8/tags"
)

// Converts a Playlist object into an m3u8 formatted string.
func EncodePlaylist(playlist *pl.Playlist) (string, error) {
	if playlist == nil || playlist.Head == nil {
		return "", fmt.Errorf("playlist is empty")
	}

	var builder strings.Builder
	builder.Grow(estimateEncodedSize(playlist))
	current := playlist.Head
	for current != nil {
		encoder, exists := tags.Encoders[current.HLSElement.Name]
		if !exists {
			return "", fmt.Errorf("unknown tag: %s", current.HLSElement.Name)
		}
		if err := encoder.Encode(current, &builder); err != nil {
			return "", fmt.Errorf("error encoding tag %s: %w", current.HLSElement.Name, err)
		}
		current = current.Next
	}
	return builder.String(), nil
}

// Estimates the encoded playlist size so the builder is allocated once, instead of growing (and copying) repeatedly.
// Each node takes roughly its URI plus a tag line; tags with many attributes (e.g. #EXT-X-STREAM-INF) take more.
func estimateEncodedSize(playlist *pl.Playlist) int {
	const tagLineSize, attrSize = 32, 16
	size := 0
	for current := playlist.Head; current != nil; current = current.Next {
		size += len(current.HLSElement.URI) + tagLineSize
		if attrs := len(current.HLSElement.Attrs); attrs > 2 {
			size += attrs * attrSize
		}
	}
	return size
}
