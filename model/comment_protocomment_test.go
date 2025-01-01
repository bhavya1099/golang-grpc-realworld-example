package model

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	pb "github.com/raahii/golang-grpc-realworld-example/proto"
)


type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type T struct {
	common
	isEnvSet bool
	context  *testContext // For running tests and subtests.
}

type Time struct {
	// wall and ext encode the wall time seconds, wall time nanoseconds,
	// and optional monotonic clock reading in nanoseconds.
	//
	// From high to low bit position, wall encodes a 1-bit flag (hasMonotonic),
	// a 33-bit seconds field, and a 30-bit wall time nanoseconds field.
	// The nanoseconds field is in the range [0, 999999999].
	// If the hasMonotonic bit is 0, then the 33-bit field must be zero
	// and the full signed 64-bit wall seconds since Jan 1 year 1 is stored in ext.
	// If the hasMonotonic bit is 1, then the 33-bit field holds a 33-bit
	// unsigned wall seconds since Jan 1 year 1885, and ext holds a
	// signed 64-bit monotonic clock reading, nanoseconds since process start.
	wall uint64
	ext  int64

	// loc specifies the Location that should be used to
	// determine the minute, hour, month, day, and year
	// that correspond to this Time.
	// The nil location means UTC.
	// All UTC times are represented with loc==nil, never loc==&utcLoc.
	loc *Location
}
func TestCommentProtoComment(t *testing.T) {

	t.Run("Return ProtoComment with correct values", func(t *testing.T) {
		comment := &Comment{
			ID:        1,
			Body:      "Test Body",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		expectedProtoComment := &pb.Comment{
			Id:        "1",
			Body:      "Test Body",
			CreatedAt: comment.CreatedAt.Format(time.RFC3339),
			UpdatedAt: comment.UpdatedAt.Format(time.RFC3339),
		}

		actualProtoComment := comment.ProtoComment()

		assert.Equal(t, expectedProtoComment, actualProtoComment, "ProtoComment should have correct values")
		t.Log("ProtoComment with correct values test passed")
	})

	t.Run("Return ProtoComment with empty values", func(t *testing.T) {
		comment := &Comment{}

		expectedProtoComment := &pb.Comment{
			Id:        "0",
			Body:      "",
			CreatedAt: "",
			UpdatedAt: "",
		}

		actualProtoComment := comment.ProtoComment()

		assert.Equal(t, expectedProtoComment, actualProtoComment, "ProtoComment should have empty values")
		t.Log("ProtoComment with empty values test passed")
	})

	t.Run("Return ProtoComment with formatted datetime", func(t *testing.T) {
		createdAt := time.Date(2022, time.January, 1, 12, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2022, time.February, 1, 12, 0, 0, 0, time.UTC)

		comment := &Comment{
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		expectedProtoComment := &pb.Comment{
			Id:        "0",
			Body:      "",
			CreatedAt: createdAt.Format(time.RFC3339),
			UpdatedAt: updatedAt.Format(time.RFC3339),
		}

		actualProtoComment := comment.ProtoComment()

		assert.Equal(t, expectedProtoComment, actualProtoComment, "ProtoComment should have formatted datetime")
		t.Log("ProtoComment with formatted datetime test passed")
	})

	t.Run("Return ProtoComment with ID as string", func(t *testing.T) {
		comment := &Comment{
			ID: 100,
		}

		expectedProtoComment := &pb.Comment{
			Id:        "100",
			Body:      "",
			CreatedAt: "",
			UpdatedAt: "",
		}

		actualProtoComment := comment.ProtoComment()

		assert.Equal(t, expectedProtoComment, actualProtoComment, "ProtoComment should have ID as string")
		t.Log("ProtoComment with ID as string test passed")
	})
}
