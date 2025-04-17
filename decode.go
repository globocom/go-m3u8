package go_m3u8

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"time"
	"unicode"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/globocom/go-m3u8/tags"
	"github.com/rs/zerolog/log"
)

type Source interface {
	io.ReadCloser
}

func ParsePlaylist(src Source) (*pl.Playlist, error) {
	playlist := &pl.Playlist{
		DoublyLinkedList: new(internal.DoublyLinkedList),
		CurrentNode:      new(internal.Node),
		CurrentSegment:   new(internal.Segment),
		CurrentStreamInf: new(internal.StreamInf),
		ProgramDateTime:  *new(time.Time),
		MediaSequence:    0,
		SegmentsCounter:  0,
		DVR:              0,
	}

	scanner := bufio.NewScanner(src)
	defer func() {
		if err := src.Close(); err != nil {
			log.Error().Err(err).Msg("error scanning playlist file")
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
			if err := pl.HandleNonTags(line, playlist); err != nil {
				return nil, fmt.Errorf("error handling non-tag line %q: %w", line, err)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to parse playlist at line: %q, error: %w", scanner.Text(), err)
	}

	return playlist, nil
}

func extractPrefix(line string) string {
	if line == "" {
		return ""
	}

	for i, r := range line {
		if r == ':' || unicode.IsSpace(r) {
			return line[:i]
		}
	}
	return line
}
