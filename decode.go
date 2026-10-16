package go_m3u8

import (
	"fmt"
	"io"
	"strings"

	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/globocom/go-m3u8/tags"
	"github.com/rs/zerolog/log"
)

// Reads an m3u8 playlist from the provided source and returns a Playlist object.
// It scans each line, identifies HLS elements, and applies the appropriate parser.
func ParsePlaylist(src io.ReadCloser) (*pl.Playlist, error) {
	playlist := pl.NewPlaylist()

	defer func() {
		if err := src.Close(); err != nil {
			log.Error().Str("service", "go-m3u8/decode.go").Err(err).Msg("error scanning playlist file")
		}
	}()

	// Reads the whole source once: every line (and so every URI and tag value stored in the playlist)
	// is a substring of this single string, instead of one allocation per line.
	var content strings.Builder
	if _, err := io.Copy(&content, src); err != nil {
		return nil, fmt.Errorf("failed to read playlist: %w", err)
	}

	remaining := content.String()
	for remaining != "" {
		var rawLine string
		rawLine, remaining, _ = strings.Cut(remaining, "\n")
		line := strings.TrimSpace(rawLine)
		linePrefix := extractPrefix(line)
		parser, exists := tags.Parsers[linePrefix]
		if exists {
			if err := parser.Parse(line, playlist); err != nil {
				return nil, fmt.Errorf("error parsing tag %s: %w", linePrefix, err)
			}
		} else {
			if err := pl.HandleMultiLineHLSElements(line, playlist); err != nil {
				return nil, fmt.Errorf("error handling multi-line HLS element %q: %w", line, err)
			}
		}
	}

	return playlist, nil
}

// Lines that start with the character '#' are either comments or tags.
// Tags begin with #EXT or #USP. All other lines that begin with '#' are comments and SHOULD be ignored.
func extractPrefix(line string) string {
	if line == "" || line[0] != '#' {
		// Not a tag or comment, return as is (URI or data line)
		return line
	}

	// Check if it's an HLS tag (#ext or #usp)
	if isRegularHLSTag(line) {
		// Extract prefix until ':' or whitespace
		idx := strings.IndexAny(line, ": \t\n\r")
		if idx != -1 {
			return line[:idx]
		}
		return line
	}

	// if starts with '#' but is not a tag, it's a comment
	return tags.CommentLineTag
}

// isRegularHLSTag checks if a line is an HLS tag by examining its prefix.
func isRegularHLSTag(line string) bool {
	return strings.HasPrefix(line, "#EXT") || strings.HasPrefix(line, "#ext") || strings.HasPrefix(line, "#USP")
}
