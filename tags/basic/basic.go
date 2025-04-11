//	Basic Tags (Section 4.4.1 in RFC)
//
// These tags are allowed in both Media Playlists and Multivariant Playlists.
package basic

import (
	"fmt"
	"strings"

	"github.com/globocom/go-m3u8/internal"
	pl "github.com/globocom/go-m3u8/playlist"
)

const (
	M3u8IdentifierTag = "#EXTM3U"
	VersionTag        = "#EXT-X-VERSION"
)

type (
	M3u8IdentifierParser struct{}
	VersionParser        struct{}
)

// Encoder Types for each tag
type (
	M3u8IdentifierEncoder struct{}
	VersionEncoder        struct{}
)

func (p M3u8IdentifierParser) Parse(tag string, playlist *pl.Playlist) error {
	playlist.Insert(&internal.Node{
		HLSElement: &internal.HLSElement{
			Name: "M3u8Identifier",
			Attrs: map[string]string{
				M3u8IdentifierTag: "",
			},
		},
	})
	return nil
}

func (p VersionParser) Parse(tag string, playlist *pl.Playlist) error {
	parts := strings.Split(tag, ":")
	if len(parts) > 1 && parts[1] != "" {
		playlist.Insert(&internal.Node{
			HLSElement: &internal.HLSElement{
				Name:  "Version",
				Attrs: map[string]string{VersionTag: strings.TrimSpace(parts[1])},
			},
		})
		return nil
	}
	return fmt.Errorf("invalid version tag: %s", tag)
}

func (e M3u8IdentifierEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	_, err := builder.WriteString(M3u8IdentifierTag + "\n")
	return err
}

func (e VersionEncoder) Encode(node *internal.Node, builder *strings.Builder) error {
	return pl.EncodeSimpleTag(node, builder, VersionTag, VersionTag)
}
