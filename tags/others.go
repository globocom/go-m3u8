//	Non-Conventional Tags
//
// The tags in this section are not traditional tags as described in the RFC.
// An example are exclusive tags added to the manifest by the packaging service.
// https://docs.unified-streaming.com/documentation/live/scte-35.html
package tags

import (
	"fmt"
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/rs/zerolog/log"
)

const (
	USPTimestampMapName = "UspTimestampMap"
	EventCueOutName     = "CueOut"
	EventCueOutContName = "CueOutCont"
	EventCueInName      = "CueIn"
	CommentLineName     = "Comment"
)

var (
	USPTimestampMapTag = "#USP-X-TIMESTAMP-MAP"
	EventCueOutTag     = "#EXT-X-CUE-OUT"
	EventCueOutContTag = "#EXT-X-CUE-OUT-CONT"
	EventCueInTag      = "#EXT-X-CUE-IN"
	CommentLineTag     = "# comment"

	USPTimestampMapOrderAttr       = []string{"MPEGTS", "LOCAL"}
	USPTimestampMapShouldQuoteAttr = map[string]bool{"MPEGTS": false, "LOCAL": false}

	EventCueOutContOrderAttr       = []string{"ElapsedTime", "Duration", "SCTE35"}
	EventCueOutContShouldQuoteAttr = map[string]bool{"ElapsedTime": false, "Duration": false, "SCTE35": false}
)

type (
	USPTimestampMapParser struct{}
	EventCueOutParser     struct{}
	EventCueOutContParser struct{}
	EventCueInParser      struct{}
	CommentParser         struct{}
)

type (
	USPTimestampMapEncoder struct{}
	EventCueOutEncoder     struct{}
	EventCueOutContEncoder struct{}
	EventCueInEncoder      struct{}
	CommentEncoder         struct{}
)

func (p USPTimestampMapParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 0 {
		params := pl.TagsToMap(parts[1])
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  USPTimestampMapName,
				Attrs: params,
			},
		})
		return nil
	}
	return fmt.Errorf("invalid usp timestamp map tag: %s", tag)
}

func (p EventCueOutParser) Parse(tag string, playlist *pl.Playlist) error {
	duration := "0"
	parts := strings.SplitN(tag, ":", 2)

	if len(parts) > 1 {
		duration = strings.TrimSpace(parts[1])
	} else {
		log.Error().Str("service", "go-m3u8/tags/others/others.go").Msgf("invalid cue out tag: %s", tag)
	}

	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  EventCueOutName,
			Attrs: map[string]string{EventCueOutTag: duration},
		},
	})

	return nil
}

func (p EventCueOutContParser) Parse(tag string, playlist *pl.Playlist) error {
	attrs := map[string]string{}
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 1 && parts[1] != "" {
		body := parts[1]
		if strings.Contains(body, "=") {
			// key=value format: ElapsedTime=5.939,Duration=201.467,SCTE35=...
			attrs = pl.TagsToMapCaseSensitive(body)
		} else {
			// positional format: 8.308/30
			if pos := strings.SplitN(body, "/", 2); len(pos) == 2 {
				attrs["ElapsedTime"] = strings.TrimSpace(pos[0])
				attrs["Duration"] = strings.TrimSpace(pos[1])
			}
		}
	}

	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  EventCueOutContName,
			Attrs: attrs,
		},
	})

	return nil
}

func (p EventCueInParser) Parse(tag string, playlist *pl.Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: EventCueInName,
			Attrs: map[string]string{
				EventCueInTag: "",
			},
		},
	})

	return nil
}

func (p CommentParser) Parse(line string, playlist *pl.Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: CommentLineName,
			Attrs: map[string]string{
				CommentLineName: line,
			},
		},
	})

	return nil
}

func (e USPTimestampMapEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeTagWithAttributes(builder, USPTimestampMapTag, node.HLSElement.Attrs, USPTimestampMapOrderAttr, USPTimestampMapShouldQuoteAttr)
}

func (e EventCueOutEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, EventCueOutTag, EventCueOutTag)
}

func (e EventCueOutContEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeTagWithAttributes(builder, EventCueOutContTag, node.HLSElement.Attrs, EventCueOutContOrderAttr, EventCueOutContShouldQuoteAttr)
}

func (e EventCueInEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(EventCueInTag + "\n")
	return err
}

func (e CommentEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	attr := fmt.Sprintf("%s\n", node.HLSElement.Attrs["Comment"])
	_, err := builder.WriteString(attr)
	return err
}
