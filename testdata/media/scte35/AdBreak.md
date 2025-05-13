# Ad Break: DATERANGE with SCTE-35 Marking

An Ad Break is present on the media playlist by way of SCTE-35 Marking. 

- The Break starts with a `EXT-X-DATERANGE` tag with attribute `SCTE-OUT` AND a `EXT-X-CUE-OUT` tag. In all cases, both tags are ALWAYS present.
- The Break ends with a `EXT-X-DATERANGE` tag with attribute `SCTE-IN` OR a `EXT-X-CUE-IN` tag. In some cases, both tags MAY be present.
- The `EXT-X-CUE-OUT` and `EXT-X-CUE-IN` tags are usually accompanied by a PDT `EXT-X-PROGRAM-DATE-TIME` tag.
- The Break's *Media Sequence* should ALWAYS equal the *Media Sequence* of the first segment inside the Break.

Depending on the playlist DVR, not all Ad Breaks will have this full structure present on the manifest at the time of request. There could be Ad Break segments that have already left the manifest due to the DVR, and Ad Breaks which are still starting and don't have the first segment on the manifest yet.

This leaves us with three possible scenarios to handle when parsing the media playlist:

## Complete Ad Breaks

Most manifests will contain Ad Breaks (one or more) that have clear start and finish markings and all Break segments present inside it.

Note that this will leave the playlist with multiple `EXT-X-PROGRAM-DATE-TIME` tags:

- The first tracks the playlist's *next media segment* following the media sequence (no Ad Break): `363987560.ts`.
- The second tracks the *next media segment* when there is an Ad Break start: `363987564.ts`.
- The third tracks the *next media segment* when there is an Ad Break end: `363987568.ts`.

Therefore, we have one PDT that tracks the playlist's media sequence, and others that track when the Ad Break starts and finishes.

**Test File:** `/testdata/media/scte35/withCompleteAdBreak.m3u8`
```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793)
#EXT-X-MEDIA-SEQUENCE:363987560
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=3873774439,LOCAL=2025-05-13T12:44:44.233300Z
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T12:44:44.233333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987560.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987561.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987562.ts
#EXTINF:5.9333, no desc
channel-audio_1=96000-video=3442944-363987563.ts
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559039-1747140304",START-DATE="2025-05-13T12:45:04.566666Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006A3F7FEFFE3775B370FE001B77400001010100001AC3CE61
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T12:45:04.566666Z
#EXTINF:3.6666, no desc
channel-audio_1=96000-video=3442944-363987564.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987565.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987566.ts
#EXTINF:6.7333, no desc
channel-audio_1=96000-video=3442944-363987567.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T12:45:24.566666Z
#EXTINF:2.8666, no desc
channel-audio_1=96000-video=3442944-363987568.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987569.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987570.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363987571.ts
```

## Partial Ad Breaks

However, there are cases where the client requests the media playlist, and the manifest will have partial Ad Breaks, either at the beginning (segments outside playlist DVR) or at the end (segments not generated yet).

These cases require special handling when parsing the media playlist information. For example, we cannot calculate the Break's *Media Sequence* if we don't have the Break's first segment present on the manifest.

### End of Playlist: Ad Break Is Not Ready Yet

It is possible that when the client requests the playlist, there will be an Ad Break start at the end of the manifest, but no segments for the Ad Break will have been generated yet. 

The `EXT-X-DATERANGE` (`SCTE-OUT`) tag is already on the manifest, as might be the `EXT-X-CUE-OUT` tag, but there are NO Break segments or PDT tag yet. 

**Test File:** `/testdata/media/scte35/withAdBreakNewNotReady.m3u8`
```
#EXTM3U
#EXT-X-VERSION:3
## Created with Unified Streaming Platform  (version=1.14.4-30793) [1eeb917c762a093f342c5a3aadcd5c8da5875e63a3ca2c61efa3bc2d8ce35227]
#EXT-X-MEDIA-SEQUENCE:363969827
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=4803017031,LOCAL=2025-05-12T13:06:05.433300Z
#EXT-X-PROGRAM-DATE-TIME:2025-05-12T13:06:05.433333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969827.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969828.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969829.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969830.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969831.ts
(...)
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969991.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969992.ts
#EXTINF:6.2333, no desc
channel-audio_1=96000-video=3442944-363969993.ts
## splice_insert()
#EXT-X-DATERANGE:ID="1-1747055968",START-DATE="2025-05-12T13:19:28.466666Z",PLANNED-DURATION=60.033333,SCTE35-OUT=0xFC3025000000000BB802FFF01405000000017FEFFF8D788E687E00527178000100000000A4C46C9A
```

Later, when the first segment for the Ad Break has been generated, we will have the `EXT-X-DATERANGE` (`SCTE-OUT`), `EXT-X-CUEOUT` and `EXT-X-PROGRAM-DATE-TIME` tags, followed by the first Break segment `EXTINF`. 

**Test File:** `/testdata/media/scte35/withAdBreakNewtReady.m3u8`
```
(...)
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969991.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363969992.ts
#EXTINF:6.2333, no desc
channel-audio_1=96000-video=3442944-363969993.ts
## splice_insert()
#EXT-X-DATERANGE:ID="1-1747055968",START-DATE="2025-05-12T13:19:28.466666Z",PLANNED-DURATION=60.033333,SCTE35-OUT=0xFC3025000000000BB802FFF01405000000017FEFFF8D788E687E00527178000100000000A4C46C9A
#EXT-X-CUE-OUT:60.033333
#EXT-X-PROGRAM-DATE-TIME:2025-05-12T13:19:28.466666Z
#EXTINF:3.3666, no desc
channel-audio_1=96000-video=3442944-363969994.ts
```

### Start of Playlist: Ad Break Leaves the DVR Limit

It is possible that when the client requests the playlist, there is an Ad Break that has already partially left the DVR range or is about to.

In the example below, the Ad Break is just about to leave the DVR, with only two non-Break segments remaining.

**Test File:** `/testdata/media/scte35/withAdBreakBeforeDVRLimit.m3u8`
```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793)
#EXT-X-MEDIA-SEQUENCE:363992684
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=6087342439,LOCAL=2025-05-13T19:34:39.433300Z
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:39.433333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992684.ts
#EXTINF:5.3666, no desc
channel-audio_1=96000-video=3442944-363992685.ts
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:49.599999Z
#EXTINF:4.2333, no desc
channel-audio_1=96000-video=3442944-363992686.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992687.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992688.ts
#EXTINF:6.1666, no desc
channel-audio_1=96000-video=3442944-363992689.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992691.ts
```
Some important changes happen to the manifest once the Ad Break reaches the DVR limit.

#### Media Sequence: Current Segment is the FIRST media segment INSIDE the Break

To avoid duplicate PDT tags, the playlist's first `EXT-X-PROGRAM-DATE` tag, which was tracking the media sequence, LEAVES the manifest, and the Break start PDT tag will take over accompaning the next media segments.

**Test File:** `/testdata/media/scte35/withAdBreakOnDVRLimit.m3u8`
```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793)
#EXT-X-MEDIA-SEQUENCE:363992686
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=6088257439,LOCAL=2025-05-13T19:34:49.599900Z
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
#EXT-X-CUE-OUT:20
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:49.599999Z
#EXTINF:4.2333, no desc
channel-audio_1=96000-video=3442944-363992686.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992687.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992688.ts
#EXTINF:6.1666, no desc
channel-audio_1=96000-video=3442944-363992689.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
```

Once the first media segment inside the Break has left the manifest, the Ad Break's *Media Sequence* will be set as **zero**.

#### Media Sequence: FIRST media segment has LEFT the Playlist

The `EXT-X-CUE-OUT` tag LEAVES the manifest alongside the first Break segment (i.e. the playlist media sequence is, at least, the SECOND media segment INSIDE the Break).

The `EXT-X-DATERANGE` (`SCTE-OUT`) tag will STAY during the Break and LEAVE only when the Break ends.

**Test File:** `/testdata/media/scte35/withAdBreakOutsideDVRLimit.m3u8`
```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793)
#EXT-X-MEDIA-SEQUENCE:363992688
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=6089070439,LOCAL=2025-05-13T19:34:58.633300Z
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:34:58.633333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992688.ts
#EXTINF:6.1666, no desc
channel-audio_1=96000-video=3442944-363992689.ts
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992691.ts
#EXTINF:4.8, no desc
```

#### Media Sequence: Current segment is the FIRST media segment OUTSIDE the Break

To avoid duplicate PDT tags, the Break start PDT `EXT-X-PROGRAM-DATE` tag, which was tracking the media sequence, LEAVES the manifest, and the Break end PDT tag will take over accompaning the next media segments.

If there is a `EXT-X-DATERANGE` tag with `SCTE-IN`, the `EXT-X-DATERANGE` (`SCTE-OUT`) tag LEAVES the manifest. Otherwise, it leaves in the next media segment.

```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793) [0db1ee867823a80c252dd4604dd94aa9d482a83e6ae968df1da76cc43321ff2d]
#EXT-X-MEDIA-SEQUENCE:363992690
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=6090057439,LOCAL=2025-05-13T19:35:09.599900Z
## splice_insert(auto_return)
#EXT-X-DATERANGE:ID="4026559475-1747164889",START-DATE="2025-05-13T19:34:49.599999Z",PLANNED-DURATION=20,SCTE35-OUT=0xFC3025000000000BB800FFF01405F0006BF37FEFFEBB581B38FE001B7740000101010000B80E326E
## Auto Return Mode
#EXT-X-CUE-IN
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:09.599999Z
#EXTINF:3.4333, no desc
channel-audio_1=96000-video=3442944-363992690.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992691.ts
#EXTINF:4.8, no desc
```

#### Media Sequence: Current segment is the SECOND media segment OUTSIDE the Break

The `EXT-X-DATERANGE` (`SCTE-OUT`) and `EXT-X-CUE-IN` tags LEAVE the manifest.

```
#EXTM3U
#EXT-X-VERSION:11
## Created with Unified Streaming Platform  (version=1.14.4-30793) [0db1ee867823a80c252dd4604dd94aa9d482a83e6ae968df1da76cc43321ff2d]
#EXT-X-MEDIA-SEQUENCE:363992691
#EXT-X-INDEPENDENT-SEGMENTS
#EXT-X-TARGETDURATION:7
#USP-X-TIMESTAMP-MAP:MPEGTS=6090366439,LOCAL=2025-05-13T19:35:13.033300Z
#EXT-X-PROGRAM-DATE-TIME:2025-05-13T19:35:13.033333Z
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992691.ts
#EXTINF:4.8, no desc
channel-audio_1=96000-video=3442944-363992692.ts
```

Now, we await until the next Ad Break, when this process will repeat.
