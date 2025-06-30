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
	ExtInfTag          = "#EXTINF"
	DiscontinuityTag   = "#EXT-X-DISCONTINUITY"
	ProgramDateTimeTag = "#EXT-X-PROGRAM-DATE-TIME"
	KeyTag             = "#EXT-X-KEY"
	ByteRangeTag       = "#EXT-X-BYTERANGE" //todo: has attributes
	MapTag             = "#EXT-X-MAP"       //todo: has attributes
	GapTag             = "#EXT-X-GAP"       //todo
	PartTag            = "#EXT-X-PART"      //todo: has attributes
)

type (
	ExtInfParser          struct{}
	DiscontinuityParser   struct{}
	ProgramDateTimeParser struct{}
	ExtKeyParser          struct{}
)

type (
	ExtInfEncoder          struct{}
	DiscontinuityEncoder   struct{}
	ProgramDateTimeEncoder struct{}
	ExtKeyEncoder          struct{}
)

func (p ExtInfParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 {
		var duration, title string

		attrs := strings.Split(parts[1], ",")
		duration = attrs[0]

		if len(attrs) > 1 {
			title = attrs[1]
		}

		playlist.CurrentSegment = pl.GetExtInfData(duration, title, playlist.MediaSequence, playlist.SegmentsCounter, playlist.DVR, playlist.ProgramDateTime)

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

func (p ExtKeyParser) Parse(tag string, playlist *pl.Playlist) error {
	params := pl.TagsToMap(tag)
	if len(params) < 1 {
		return fmt.Errorf("invalid ext key tag: %s", tag)
	}

	// METHOD attribute is REQUIRED by RFC
	if params["METHOD"] == "" {
		return fmt.Errorf("METHOD attribute is required: %s", tag)
	}

	// URI attribute is REQUIRED unless the METHOD is NONE
	if (params["METHOD"] != "NONE") && (params["URI"] == "") {
		return fmt.Errorf("URI attribute is required when METHOD is not NONE: %s", tag)
	}

	// IV attribute is required if METHOD is AES-128
	if (params["METHOD"] == "AES-128") && params["IV"] == "" {
		return fmt.Errorf("IV attribute is required when METHOD is AES-128: %s", tag)
	}

	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "ExtKey",
			Attrs: params,
		},
	})

	return nil
}

func (e ExtInfEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	duration := node.HLSElement.Attrs["Duration"]
	title := node.HLSElement.Attrs["Title"]
	uri := node.HLSElement.URI

	// #EXTINF:<duration>,[<title>]
	if title != "" {
		title = "," + title
	}

	attr := fmt.Sprintf("%s:%s%s\n%s\n", ExtInfTag, duration, title, uri)
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

func (e ExtKeyEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	orderAttr := []string{"METHOD", "URI", "IV", "KEYFORMAT", "KEYFORMATVERSIONS"}
	shouldQuoteAttr := map[string]bool{
		"METHOD":            false,
		"URI":               true,
		"IV":                false,
		"KEYFORMAT":         true,
		"KEYFORMATVERSIONS": true,
	}
	return pl.EncodeTagWithAttributes(builder, KeyTag, node.HLSElement.Attrs, orderAttr, shouldQuoteAttr)
}
