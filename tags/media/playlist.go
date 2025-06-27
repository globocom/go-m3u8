//	Media Playlist Tags (Section 4.4.3 on RFC)
//
// Media Playlist tags describe global parameters of the Media Playlist.
// There MUST NOT be more than one Media Playlist tag of each type in
// any Media Playlist.
//
// A Media Playlist tag MUST NOT appear in a Multivariant Playlist.
package media

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
)

var (
	TargetDurationTag        = "#EXT-X-TARGETDURATION"
	MediaSequenceTag         = "#EXT-X-MEDIA-SEQUENCE"
	DiscontinuitySequenceTag = "#EXT-X-DISCONTINUITY-SEQUENCE" //todo: has one attribute
	EndlistTag               = "#EXT-X-ENDLIST"                //todo
	PlaylistTypeTag          = "#EXT-X-PLAYLIST-TYPE"          //todo: has one attribute
	IFramesOnlyTag           = "#EXT-X-I-FRAMES-ONLY"          //todo
	PartInfTag               = "#EXT-X-PART-INF"               //todo: has attributes
	ServerControlTag         = "#EXT-X-SERVER-CONTROL"         //todo: has attributes
)

type (
	TargetDurationParser struct{}
	MediaSequenceParser  struct{}
)

type (
	TargetDurationEncoder struct{}
	MediaSequenceEncoder  struct{}
)

func (p TargetDurationParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "TargetDuration",
				Attrs: map[string]string{TargetDurationTag: strings.TrimSpace(parts[1])},
			},
		})
		return nil
	}
	return fmt.Errorf("invalid target duration tag: %s", tag)
}

func (p MediaSequenceParser) Parse(tag string, playlist *pl.Playlist) error {
	var err error
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		mediaSequence := strings.TrimSpace(parts[1])
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "MediaSequence",
				Attrs: map[string]string{MediaSequenceTag: mediaSequence},
			},
		})

		playlist.MediaSequence, err = strconv.Atoi(mediaSequence)
		if err != nil {
			return fmt.Errorf("invalid media sequence number: %s", tag)
		}

		return nil
	}
	return fmt.Errorf("invalid media sequence tag: %s", tag)
}

func (e TargetDurationEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, TargetDurationTag, TargetDurationTag)
}

func (e MediaSequenceEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, MediaSequenceTag, MediaSequenceTag)
}
