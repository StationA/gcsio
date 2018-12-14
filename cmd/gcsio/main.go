package main

import (
	"context"
	"io"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gobwas/glob"
	"google.golang.org/api/iterator"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	uploadCmd    = kingpin.Command("upload", "Streams stdin up to a GCS object")
	dstUri       = uploadCmd.Arg("dst", "Destination object URI").Required().URL()
	catCmd       = kingpin.Command("cat", "Streams an object from GCS to stdout")
	noDecompress = catCmd.Flag("no-decompress", "Disable automatic stream decompression").Bool()
	srcUri       = catCmd.Arg("source", "Source object URI or glob pattern").Required().URL()
)

func doUpload(ctx context.Context) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	f := os.Stdin

	objectUri := *dstUri
	wc := client.Bucket(objectUri.Host).Object(objectUri.Path[1:]).NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func getMatchingObjects(ctx context.Context, bucket *storage.BucketHandle, pattern string) ([]string, error) {
	var matchingObjects []string

	p := glob.MustCompile(pattern)
	it := bucket.Objects(ctx, nil)
	for {
		obj, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		if p.Match(obj.Name) {
			matchingObjects = append(matchingObjects, obj.Name)
		}
	}
	return matchingObjects, nil
}

func doStream(ctx context.Context) error {
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v\n", err)
	}

	f := os.Stdout

	objectUri := *srcUri
	bucket := client.Bucket(objectUri.Host)
	pattern := objectUri.Path[1:]

	matchingObjects, err := getMatchingObjects(ctx, bucket, pattern)
	if err != nil {
		return err
	}
	multireader := NewMultiObjectReader(ctx, bucket, matchingObjects)

	if _, err = io.Copy(f, multireader); err != nil {
		return err
	}
	return nil
}

func main() {
	ctx := context.Background()
	cmd := kingpin.Parse()
	var err error
	switch cmd {
	case "upload":
		err = doUpload(ctx)
	case "cat":
		err = doStream(ctx)
	}

	if err != nil {
		panic(err)
	}
}
