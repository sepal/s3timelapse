package main

import "testing"

func TestURLParsing(t *testing.T) {
	bucket, prefix := ParseUrl(("s3://mybucket/some/path"))

	if bucket != "mybucket" {
		t.Fatalf("Bucket name %s doesn't matches expected name 'mybucket'", bucket)
	}

	if prefix != "some/path" {
		t.Fatalf("Prefix %s doesn't matches 'some/path'", prefix)
	}
}
