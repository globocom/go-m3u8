<div align=center><img src="gopher.png" height=100px></div>
<h1 align=center>
go-m3u8
<div align=center>
<img src="https://github.com/globocom/go-m3u8/actions/workflows/go.yml/badge.svg">
<a href="https://goreportcard.com/report/github.com/globocom/go-m3u8"><img src="https://goreportcard.com/badge/github.com/globocom/go-m3u8"/></a>
<img src="https://img.shields.io/github/go-mod/go-version/globocom/go-m3u8">
</div>
</h1>


### Work in Progress!

_There isn't a stable version for now. As this is currently a WIP, the API may have changes._

## go-m3u8

Parser for [m3u8](https://tools.ietf.org/html/rfc8216) to facilitate manifest manipulation.

_We currently only support Live Streaming manifests._

### 1. Doubly Linked List

The m3u8 is represented by a doubly linked list. This data structure allows us to access the manifest in a sorted manner and apply operations (modify, add, remove) to its content.

Some available operations are:

- **Create** a new node to represent an HLS element of the m3u8 manifest (e.g. a tag, a comment).
- **Insert** a new node into the playlist (append at the end OR insert before/after/between specific nodes).
- **Find** a specific node or a list of nodes based on the HLS element name.

### 2. Tags

To guarantee scalibility, our lib considers a data structure for the HLS elements that follows the [RFC documentation](https://tools.ietf.org/html/rfc8216).

The [**tags**](https://github.com/globocom/go-m3u8/tags) package separates HLS elements into sub-packages according to the [**Playlist Tags**](https://datatracker.ietf.org/doc/html/draft-pantos-hls-rfc8216bis#section-4.4) section on the RFC.

(The tags listed below are the ones currently supported by the lib.)

1. **basic -** Basic Tags (Section 4.4.1).
- `#EXTM3U`
- `#EXT-X-VERSION`

2. **exclusive -** Media or Multivariant Playlist Tags (Section 4.4.2).
- `#EXT-X-INDEPENDENT-SEGMENTS`
- `#EXT-X-DEFINE`

3. **media -** Media Playlist, Metadata and Segment Tags (Sections 4.4.3 to 4.4.5).
- `#EXT-X-DATERANGE`
- `#EXT-X-TARGETDURATION`
- `#EXT-X-MEDIA-SEQUENCE`
- `#EXT-X-DISCONTINUITY-SEQUENCE`
- `#EXTINF`
- `#EXT-X-DISCONTINUITY`
- `#EXT-X-PROGRAM-DATE-TIME`
- `#EXT-X-KEY`

4. **multivariant -** Multivariant Playlist Tags (Section 4.4.6).
- `#EXT-X-STREAM-INF`

5. **others -** The tags in this section are "non-official" and are not listed in the RFC, e.g. tags added to the manifest by the packaging service.
- `#EXT-X-CUE-OUT`
- `#EXT-X-CUE-IN`
- Packager specific tags.
- In-manifest comments (begin with `#` and are NOT tags).

## Getting started

```
go install github.com/globocom/go-m3u8
```