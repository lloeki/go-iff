// Copyright 2014 Loic Nageleisen. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package iff provides interfaces to the Interchange File Format
*/

package iff

import (
    "os"
    "fmt"
    "io"
    "errors"
    "encoding/binary"
)

const (
    FORM = "FORM"
    LIST = "LIST"
    CAT  = "CAT"
    PROP = "PROP"
)

type Id [4]byte

func (id Id) String() string {
    return string(id[:])
}

// Chunk is the basic unit of IFF files
type Chunk struct {
    id Id
    size uint32
    data *io.SectionReader
}

// FormatChunk is a FORM chunk
type FormChunk struct {
    typ Id
    data *io.SectionReader
}

// ChunkReader is the interface of a Reader able to read chunks
type ChunkReader interface {
    io.Reader
    io.ReaderAt
    io.Seeker
}

func ReadChunk(r ChunkReader) (chunk *Chunk, err error) {
    chunk = &Chunk{}

    _, err = r.Read(chunk.id[:])
    if err != nil {
        return
    }

    err = binary.Read(r, binary.BigEndian, &chunk.size)
    if err != nil {
        return
    }

    offset, _ := r.Seek(0, os.SEEK_CUR)
    chunk.data = io.NewSectionReader(r, offset, int64(chunk.size))

    return
}

func ReadFormChunk(r ChunkReader) (chunk *FormChunk, err error) {
    c, err := ReadChunk(r)
    if err != nil {
        return
    }

    if c.id.String() != FORM {
        msg := fmt.Sprintf("expected FORM, got %s", c.id)
        err = errors.New(msg)
        return
    }

    chunk = &FormChunk{}

    _, err = c.data.Read(chunk.typ[:])
    if err != nil {
        return
    }

    offset, _ := c.data.Seek(0, os.SEEK_CUR)
    size := c.size - uint32(len(chunk.typ))
    chunk.data = io.NewSectionReader(c.data, offset, int64(size))

    return
}
