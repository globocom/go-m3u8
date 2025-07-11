//	Non-Conventional Tags
//
// The tags in this section are not traditional tags as described in the RFC.
// An example are exclusive tags added to the manifest by the packaging service.
// https://docs.unified-streaming.com/documentation/live/scte-35.html
package others

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/rs/zerolog/log"
)

var (
	USPTimestampMapTag = "#USP-X-TIMESTAMP-MAP"
	EventCueOutTag     = "#EXT-X-CUE-OUT"
	EventCueInTag      = "#EXT-X-CUE-IN"
	CommentLineTag     = "# comment"
	CommentLineRegex   = regexp2.MustCompile(`^#(?!(EXT|ext|USP)).*`, 0) // excludes tags (#EXT, #ext or #USP)
)

type (
	USPTimestampMapParser struct{}
	EventCueOutParser     struct{}
	EventCueInParser      struct{}
	CommentParser         struct{}
)

type (
	USPTimestampMapEncoder struct{}
	EventCueOutEncoder     struct{}
	EventCueInEncoder      struct{}
	CommentEncoder         struct{}
)

func (p USPTimestampMapParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.SplitN(tag, ":", 2)
	if len(parts) > 0 {
		params := pl.TagsToMap(parts[1])
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "UspTimestampMap",
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
		log.Error().Msgf("invalid cue out tag: %s", tag)
	}

	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name:  "CueOut",
			Attrs: map[string]string{EventCueOutTag: duration},
		},
	})

	return nil
}

func (p EventCueInParser) Parse(tag string, playlist *pl.Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "CueIn",
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
			Name: "Comment",
			Attrs: map[string]string{
				"Comment": line,
			},
		},
	})

	return nil
}

func (e USPTimestampMapEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	order := []string{"MPEGTS", "LOCAL"}
	return pl.EncodeTagWithAttributes(builder, USPTimestampMapTag, node.HLSElement.Attrs, order)
}

func (e EventCueOutEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, EventCueOutTag, EventCueOutTag)
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
