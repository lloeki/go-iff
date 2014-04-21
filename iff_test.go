package iff

import (
    "testing"
    "os"
    "io/ioutil"
    check "gopkg.in/check.v1"
)

func Test(t *testing.T) { check.TestingT(t) }

type IffSuite struct{}

var _ = check.Suite(&IffSuite{})

func (s *IffSuite) TestIdToString(c *check.C) {
    id := Id { 'f', 'o', 'o', 'b' }

    c.Assert(id.String(), check.Equals, "foob")
}

func (s *IffSuite) TestIdLength(c *check.C) {
    id := Id { 'f', 'o', 'o', 'b' }

    c.Assert(len(id), check.Equals, 4)
}

func (s *IffSuite) TestReadChunk(c *check.C) {
    f, err := os.Open("fixture1.iff")
    if err != nil {
        panic(err)
    }

    defer f.Close()

    chunk, err := ReadChunk(f)
    if err != nil {
        panic(err)
    }

    c.Assert(chunk.id, check.Equals, Id{'F', 'O', 'R', 'M'})
    c.Assert(chunk.size, check.Equals, uint32(26))
    c.Assert(chunk.data, check.NotNil)
}

func (s *IffSuite) TestReadFormChunk(c *check.C) {
    f, err := os.Open("fixture1.iff")
    if err != nil {
        panic(err)
    }

    defer f.Close()

    chunk, err := ReadFormChunk(f)
    if err != nil {
        panic(err)
    }

    c.Assert(chunk.typ, check.Equals, Id{'F', 'A', 'K', 'E'})
    c.Assert(chunk.data, check.NotNil)

    data_chunk, err := ReadChunk(chunk.data)
    if err != nil {
        panic(err)
    }

    c.Assert(data_chunk.id, check.Equals, Id{'D', 'A', 'T', 'A'})
    c.Assert(data_chunk.size, check.Equals, uint32(14))
    c.Assert(data_chunk.data, check.NotNil)

    data, err := ioutil.ReadAll(data_chunk.data)
    if err != nil {
        panic(err)
    }

    c.Assert(string(data), check.Equals, "Hello, world!\n")
}
