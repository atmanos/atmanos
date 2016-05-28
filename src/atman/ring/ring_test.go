package ring

import (
	"bytes"
	"testing"
)

func TestBufferWrite(t *testing.T) {
	for i, tt := range []struct {
		readerPos, writerPos uint32

		data []byte

		expectedSent int
		expectedData []byte
	}{
		{
			readerPos: 0,
			writerPos: 0,

			data: []byte("****"),

			expectedSent: 4,
			expectedData: []byte("****...."),
		},
		{
			readerPos: 0,
			writerPos: 0,

			data: []byte("*********"),

			expectedSent: 8,
			expectedData: []byte("********"),
		},
		{
			readerPos: 0,
			writerPos: 0,

			data: []byte("******************"),

			expectedSent: 8,
			expectedData: []byte("********"),
		},
		{
			readerPos: 6,
			writerPos: 6,

			data: []byte("****"),

			expectedSent: 4,
			expectedData: []byte("**....**"),
		},
		{
			readerPos: 7,
			writerPos: 13,

			data: []byte("****"),

			expectedSent: 2,
			expectedData: []byte(".....**."),
		},
		{
			readerPos: 0xFFFFFFFF - 1,
			writerPos: 0xFFFFFFFF - 1,

			data: []byte("****"),

			expectedSent: 4,
			expectedData: []byte("**....**"),
		},
	} {
		var (
			readerPos = tt.readerPos
			writerPos = tt.writerPos

			buf = &Buffer{
				Data:      []byte("........"),
				ReaderPos: &readerPos,
				WriterPos: &writerPos,
			}
		)

		sent := buf.Write(tt.data)

		if sent != tt.expectedSent {
			t.Errorf("%d. sent = %d, want %d", i, sent, tt.expectedSent)
		}

		if !bytes.Equal(buf.Data, tt.expectedData) {
			t.Errorf("%d. data = %s, want %s", i, buf.Data, tt.expectedData)
		}

		if readerPos != tt.readerPos {
			t.Errorf("%d. readerPos = %d, want %d", i, readerPos, tt.readerPos)
		}

		exptectedWriterPos := tt.writerPos + uint32(sent)
		if writerPos != exptectedWriterPos {
			t.Errorf("%d. writerPos = %d, want %d", i, writerPos, exptectedWriterPos)
		}
	}
}

func TestBufferRead(t *testing.T) {
	for i, tt := range []struct {
		readerPos, writerPos uint32

		data []byte

		expectedRead int
		expectedData []byte
	}{
		{
			readerPos: 0,
			writerPos: 0,

			data: []byte("****"),

			expectedRead: 0,
			expectedData: []byte{},
		},
		{
			readerPos: 0,
			writerPos: 8,

			data: []byte("*********"),

			expectedRead: 8,
			expectedData: []byte("********"),
		},
		{
			readerPos: 6,
			writerPos: 10,

			data: []byte("**....**"),

			expectedRead: 4,
			expectedData: []byte("****"),
		},
		{
			readerPos: 0xFFFFFFFF - 1,
			writerPos: 2,

			data: []byte("**....**"),

			expectedRead: 4,
			expectedData: []byte("****"),
		},
	} {
		var (
			readerPos = tt.readerPos
			writerPos = tt.writerPos

			buf = &Buffer{
				Data:      tt.data,
				ReaderPos: &readerPos,
				WriterPos: &writerPos,
			}
		)

		data := make([]byte, 64)
		read := buf.Read(data)

		if read != tt.expectedRead {
			t.Errorf("%d. read = %d, want %d", i, read, tt.expectedRead)
		}

		data = data[:read]

		if !bytes.Equal(data, tt.expectedData) {
			t.Errorf("%d. data = %s, want %s", i, data, tt.expectedData)
		}

		if writerPos != tt.writerPos {
			t.Errorf("%d. writerPos = %d, want %d", i, writerPos, tt.writerPos)
		}

		exptectedReaderPos := tt.readerPos + uint32(read)
		if readerPos != exptectedReaderPos {
			t.Errorf("%d. readerPos = %d, want %d", i, readerPos, exptectedReaderPos)
		}
	}
}
