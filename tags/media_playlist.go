//	Media Playlist Tags (Section 4.4.3 on RFC)
//
// Media Playlist tags describe global parameters of the Media Playlist.
// There MUST NOT be more than one Media Playlist tag of each type in
// any Media Playlist.
//
// A Media Playlist tag MUST NOT appear in a Multivariant Playlist.
package tags

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
	"github.com/rs/zerolog/log"
)

const (
	TargetDurationName        = "TargetDuration"
	MediaSequenceName         = "MediaSequence"
	DiscontinuitySequenceName = "DiscontinuitySequence"
	IFramesOnlyName           = "IFramesOnly"
)

var (
	TargetDurationTag        = "#EXT-X-TARGETDURATION"
	MediaSequenceTag         = "#EXT-X-MEDIA-SEQUENCE"
	DiscontinuitySequenceTag = "#EXT-X-DISCONTINUITY-SEQUENCE"
	IFramesOnlyTag           = "#EXT-X-I-FRAMES-ONLY"
	EndlistTag               = "#EXT-X-ENDLIST"        // todo
	PlaylistTypeTag          = "#EXT-X-PLAYLIST-TYPE"  // todo: has one attribute
	PartInfTag               = "#EXT-X-PART-INF"       // todo: has attributes
	ServerControlTag         = "#EXT-X-SERVER-CONTROL" // todo: has attributes
)

type (
	TargetDurationParser        struct{}
	MediaSequenceParser         struct{}
	DiscontinuitySequenceParser struct{}
	IFramesOnlyParser           struct{}
)

type (
	TargetDurationEncoder        struct{}
	MediaSequenceEncoder         struct{}
	DiscontinuitySequenceEncoder struct{}
	IFramesOnlyEncoder           struct{}
)

func (p TargetDurationParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  TargetDurationName,
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
				Name:  MediaSequenceName,
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

// #EXT-X-DISCONTINUITY-SEQUENCE:<number>
func (p DiscontinuitySequenceParser) Parse(tag string, playlist *pl.Playlist) error {
	var err error
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		discontinuitySequence := strings.TrimSpace(parts[1])
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  DiscontinuitySequenceName,
				Attrs: map[string]string{DiscontinuitySequenceTag: discontinuitySequence},
			},
		})

		playlist.DiscontinuitySequence, err = strconv.Atoi(discontinuitySequence)
		if err != nil {
			return fmt.Errorf("invalid discontinuity sequence number: %s", tag)
		}

		return nil
	}
	return fmt.Errorf("invalid discontinuity sequence tag: %s", tag)
}

func (p IFramesOnlyParser) Parse(tag string, playlist *pl.Playlist) error {
	hlsVersion, err := strconv.Atoi(playlist.VersionValue())
	if err != nil {
		log.Warn().Str("service", "go-m3u8/tags/media_playlist.go").Msg("could not parse playlist version")
	}

	if hlsVersion < 4 {
		return fmt.Errorf("use of %s REQUIRES a compatibility version number of 4 or greater, but playlist version is %d", tag, hlsVersion)
	}

	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: IFramesOnlyName,
			Attrs: map[string]string{
				IFramesOnlyTag: "",
			},
		},
	})
	return nil
}

func (e TargetDurationEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, TargetDurationTag, TargetDurationTag)
}

func (e MediaSequenceEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, MediaSequenceTag, MediaSequenceTag)
}

func (e DiscontinuitySequenceEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, DiscontinuitySequenceTag, DiscontinuitySequenceTag)
}

func (e IFramesOnlyEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(IFramesOnlyTag + "\n")
	return err
}
