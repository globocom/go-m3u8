package go_m3u8_test

import (
	"io"
	"os"
	"strings"
	"testing"

	m3u8 "github.com/globocom/go-m3u8"
)

func loadFixture(b *testing.B, path string) string {
	b.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		b.Fatalf("failed to read fixture %s: %v", path, err)
	}
	return string(data)
}

func benchmarkParse(b *testing.B, path string) {
	manifest := loadFixture(b, path)
	b.SetBytes(int64(len(manifest)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := m3u8.ParsePlaylist(io.NopCloser(strings.NewReader(manifest))); err != nil {
			b.Fatal(err)
		}
	}
}

func benchmarkEncode(b *testing.B, path string) {
	manifest := loadFixture(b, path)
	p, err := m3u8.ParsePlaylist(io.NopCloser(strings.NewReader(manifest)))
	if err != nil {
		b.Fatal(err)
	}
	b.SetBytes(int64(len(manifest)))
	b.ReportAllocs()
	for b.Loop() {
		if _, err := m3u8.EncodePlaylist(p); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseMedia(b *testing.B) {
	benchmarkParse(b, "testlocal/media.m3u8")
}

func BenchmarkParseMultivariant(b *testing.B) {
	benchmarkParse(b, "testlocal/multivariant.m3u8")
}

func BenchmarkEncodeMedia(b *testing.B) {
	benchmarkEncode(b, "testlocal/media.m3u8")
}

func BenchmarkEncodeMultivariant(b *testing.B) {
	benchmarkEncode(b, "testlocal/multivariant.m3u8")
}
