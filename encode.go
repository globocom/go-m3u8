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
