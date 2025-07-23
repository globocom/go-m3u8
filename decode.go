package go_m3u8

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"

	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/globocom/go-m3u8/tags"
	"github.com/globocom/go-m3u8/tags/others"
	"github.com/rs/zerolog/log"
)

type Source interface {
	io.ReadCloser
}

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
// Tags begin with #EXT.  They are case sensitive.  All other lines that begin with '#' are comments and SHOULD be ignored.
func extractPrefix(line string) string {
	// check for blank lines
	if line == "" {
		return ""
	}

	// check for comments
	isComment, err := others.CommentLineRegex.MatchString(line)
	if err != nil {
		log.Error().Str("service", "go-m3u8/decode.go").Err(err).Msgf("failed to parse line: %s", line)
		return ""
	}

	if isComment {
		return others.CommentLineTag
	}

	// check for tags and uri
	for i, r := range line {
		if r == ':' || unicode.IsSpace(r) {
			return line[:i]
		}
	}
	return line
}
