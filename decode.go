package go_m3u8

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/globocom/go-m3u8/tags"
	"github.com/rs/zerolog/log"
)

type Source interface {
	io.ReadCloser
}

// Reads an m3u8 playlist from the provided source and returns a Playlist object.
// It scans each line, identifies HLS elements, and applies the appropriate parser.
func ParsePlaylist(src Source) (*pl.Playlist, error) {
	playlist := pl.NewPlaylist()

	scanner := bufio.NewScanner(src)
	defer func() {
		if err := src.Close(); err != nil {
			log.Error().Str("service", "go-m3u8/decode.go").Err(err).Msg("error scanning playlist file")
		}
	}()

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
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
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse playlist at line: %q, error: %w", scanner.Text(), err)
	}

	return playlist, nil
}

// Lines that start with the character '#' are either comments or tags.
// Tags begin with #EXT or #USP. All other lines that begin with '#' are comments and SHOULD be ignored.
func extractPrefix(line string) string {
	if line == "" {
		return ""
	}

	// Check if line is a tag
	if isRegularHLSTag(line) {
		// it's a tag, extract prefix until ':' or whitespace
		for i, r := range line {
			if r == ':' || unicode.IsSpace(r) {
				return line[:i]
			}
		}
		return line
	}

	// if starts with '#' but is not a tag, it's a comment
	if strings.HasPrefix(line, "#") {
		return tags.CommentLineTag
	}

	// otherwise, return as is (URI or data line)
	return line
}

// isHLSTag checks if a line is an HLS tag by examining its prefix.
// Tags must start with #ext (case-insensitive) or #usp (case-insensitive).
// This function converts the prefix to lowercase internally for comparison.
func isRegularHLSTag(line string) bool {
	if len(line) < 4 {
		return false
	}

	prefix := line[:4]
	lowerPrefix := strings.ToLower(prefix)

	// Check for #ext or #usp (case-insensitive)
	if lowerPrefix == "#ext" || lowerPrefix == "#usp" {
		return true
	}

	return false
}
