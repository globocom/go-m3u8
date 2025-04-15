//	Media Segment Tags (Section 4.4.4 in RFC)
//
// Each Media Segment is specified by a series of Media Segment tags
// followed by a URI.  Some Media Segment tags apply to just the next
// segment; others apply to all subsequent segments until another
// instance of the same tag.
//
// A Media Segment tag MUST NOT appear in a Multivariant Playlist.
// Clients MUST fail to parse Playlists that contain both Media Segment
// tags and Multivariant Playlist tags
package media

import (
	"fmt"
	"strings"
	"time"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
)

var (
	ExtInfTag          = "#EXTINF"          //todo: has attributes
	ByteRangeTag       = "#EXT-X-BYTERANGE" //todo: has attributes
	DiscontinuityTag   = "#EXT-X-DISCONTINUITY"
	KeyTag             = "#EXT-X-KEY"               //todo: has attributes
	MapTag             = "#EXT-X-MAP"               //todo: has attributes
	ProgramDateTimeTag = "#EXT-X-PROGRAM-DATE-TIME" //todo: has attributes
	GapTag             = "#EXT-X-GAP"
	PartTag            = "#EXT-X-PART" //todo: has attributes
)

type (
	ExtInfParser          struct{}
	DiscontinuityParser   struct{}
	ProgramDateTimeParser struct{}
)

type (
	ExtInfEncoder          struct{}
	DiscontinuityEncoder   struct{}
	ProgramDateTimeEncoder struct{}
)

func (p ExtInfParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 {
		duration := strings.TrimSpace(strings.Split(parts[1], ",")[0])

		playlist.CurrentSegment = internal.ToSegmentType(duration, playlist.MediaSequence, playlist.SegmentsCounter, playlist.DVR, playlist.ProgramDateTime)

		playlist.DVR = pl.RoundFloat(playlist.DVR+playlist.CurrentSegment.Duration, 4)
		playlist.SegmentsCounter += 1

		return nil
	}
	return fmt.Errorf("invalid extension tag: %s", tag)
}

func (p DiscontinuityParser) Parse(tag string, playlist *pl.Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "Discontinuity",
			Attrs: map[string]string{
				DiscontinuityTag: "",
			},
		},
	})
	return nil
}

func (p ProgramDateTimeParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.SplitN(tag, ":", 2)

	if len(parts) <= 1 {
		return fmt.Errorf("invalid program date time tag: %s", tag)
	}

	dateTimeValue := strings.TrimSpace(parts[1])
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "ProgramDateTime",
			Attrs: map[string]string{ProgramDateTimeTag: dateTimeValue},
		},
	})

	if playlist.ProgramDateTime.Format(time.DateOnly) == "0001-01-01" {
		parsedTime, err := time.Parse(time.RFC3339Nano, dateTimeValue)

		if err != nil {
			return fmt.Errorf("invalid program date time tag: %s", tag)
		}

		playlist.ProgramDateTime = parsedTime
	}

	return nil
}

func (e ExtInfEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	attr := fmt.Sprintf("%s:%s\n%s\n", ExtInfTag, node.HLSElement.Attrs["Duration"], node.HLSElement.URI)
	_, err := builder.WriteString(attr)
	return err
}

func (e DiscontinuityEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(DiscontinuityTag + "\n")
	return err
}

func (e ProgramDateTimeEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, ProgramDateTimeTag, ProgramDateTimeTag)
}
