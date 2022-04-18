package main

import (
	"testing"
	"time"
)

func TestURLParsing(t *testing.T) {
	bucket, prefix := ParseUrl(("s3://mybucket/some/path"))

	if bucket != "mybucket" {
		t.Fatalf("Bucket name %s doesn't matches expected name 'mybucket'", bucket)
	}

	if prefix != "some/path" {
		t.Fatalf("Prefix %s doesn't matches 'some/path'", prefix)
	}
}

func TestIsDateInDay(t *testing.T) {
	day, _ := time.Parse("2006-01-02", "2022-04-18")

	tm, _ := time.Parse("2006-01-02 15:04", "2022-04-18 00:00")

	if !TimeInDay(tm, day) {
		t.Fatalf("Time %s should be in day!", tm)
	}

	tm, _ = time.Parse("2006-01-02 15:04", "2022-04-19 00:00")

	if TimeInDay(tm, day) {
		t.Fatalf("Time %s should be in day!", tm)
	}
}
