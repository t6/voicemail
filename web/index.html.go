package web

import (
	"bytes"
	"compress/gzip"
	"io"
)

// index_html returns the raw, uncompressed file data data.
func index_html() []byte {
	gz, err := gzip.NewReader(bytes.NewBuffer([]byte{
0x1f,0x8b,0x08,0x00,0x00,0x09,0x6e,0x88,0x00,0xff,0x94,0x57,
0xdb,0x6e,0xdb,0x38,0x13,0xbe,0xcf,0x53,0xb0,0xfc,0x81,0x5e,
0xfd,0xb2,0x9a,0x9e,0xeb,0xd8,0x06,0x8a,0xb6,0xc0,0x16,0xe8,
0x1e,0x2e,0x7a,0xb3,0x28,0x0a,0x83,0x16,0x69,0x8b,0x29,0x45,
0xaa,0xe4,0xc8,0x71,0x76,0xb1,0x77,0xfb,0x66,0xfb,0x62,0x3b,
0x3c,0x48,0xa6,0x6c,0x6d,0x92,0x06,0x88,0x45,0xcd,0x0c,0xe7,
0xf8,0xcd,0x90,0x5a,0x3c,0x7a,0xff,0xeb,0xbb,0xcf,0xbf,0xff,
0xf6,0x81,0xd4,0xd0,0xa8,0xd5,0xc5,0xc2,0x3f,0x88,0x62,0x7a,
0xb7,0xa4,0x5c,0xd0,0xd5,0x05,0x21,0x8b,0x5a,0x30,0xee,0x17,
0xb8,0x6c,0x04,0x30,0x52,0xd5,0xcc,0x3a,0x01,0x4b,0xda,0xc1,
0xb6,0x78,0x4d,0x73,0x56,0x0d,0xd0,0x16,0xe2,0x7b,0x27,0xf7,
0x4b,0xfa,0xce,0x68,0x10,0x1a,0x8a,0xcf,0xb7,0xad,0xa0,0xa4,
0x8a,0x6f,0x4b,0x0a,0xe2,0x00,0xa5,0xb7,0x72,0x35,0x28,0x1a,
0xe9,0x01,0x09,0x4a,0xac,0xde,0x6a,0xdb,0x6d,0x37,0x82,0x69,
0xb8,0x31,0x16,0x84,0x5d,0x94,0x91,0x9e,0xd9,0xd2,0xac,0x11,
0xde,0x49,0x57,0x59,0xd9,0x82,0x34,0x3a,0x33,0x42,0xcf,0x05,
0x59,0x07,0xb5,0xb1,0x63,0x99,0x28,0xf4,0xa8,0x28,0xc8,0x27,
0x41,0x7e,0xfa,0xfc,0xf3,0xa7,0x17,0xc4,0xd5,0xb2,0xf9,0x3f,
0xd9,0x1a,0x4b,0x3e,0x7e,0x78,0x59,0xbc,0x26,0xae,0x6b,0x5b,
0xf4,0x80,0x98,0x6d,0x10,0x20,0x42,0x89,0x06,0xb7,0x3b,0x52,
0x14,0xab,0x61,0xfb,0x17,0xb9,0x25,0x0a,0x70,0x07,0x79,0xf3,
0x35,0x52,0x03,0x27,0x7a,0x46,0x9c,0xad,0x96,0xd4,0x67,0x66,
0x5e,0x86,0xc0,0x5f,0x78,0x1b,0xb3,0x9d,0x31,0x3b,0x25,0x2a,
0xc3,0xc5,0xac,0x32,0x4d,0xe9,0xf6,0xba,0x04,0xdb,0xe9,0x6f,
0x51,0x64,0x76,0xed,0xe8,0x6a,0x51,0x46,0x0d,0x99,0xca,0x47,
0x5f,0x84,0xe6,0x72,0xfb,0xd5,0x5b,0x1f,0x79,0xef,0xe0,0x56,
0x89,0xcc,0x2b,0x25,0xf5,0x37,0x52,0x5b,0xb1,0x5d,0xd2,0x2d,
0xdb,0x4b,0x8c,0x7a,0x86,0x3f,0x94,0x58,0xa1,0x96,0xd4,0x61,
0x26,0xa0,0xea,0x80,0x78,0x3a,0x25,0x80,0x15,0x5a,0x52,0xd9,
0xb0,0x9d,0x28,0x0f,0x45,0xa0,0x9d,0x6b,0x49,0x01,0x54,0x9d,
0xb5,0x18,0xff,0x6c,0x63,0x0c,0x38,0xb0,0xac,0xad,0xb8,0x0e,
0x01,0x0c,0x84,0x62,0xff,0xf4,0xc9,0xf3,0xb2,0x72,0x2e,0x23,
0xa1,0xc0,0x46,0x6a,0xc1,0x67,0x8d,0x44,0x69,0xe7,0x7a,0x3f,
0x82,0xd3,0xb5,0x10,0xd0,0x1b,0x0c,0x94,0xe4,0x50,0x40,0x8a,
0x17,0xee,0xe3,0xdf,0x18,0x7e,0x4b,0xfe,0x1c,0x92,0xd1,0x32,
0xce,0xa5,0xde,0x15,0x60,0xda,0x39,0x79,0xf9,0xa4,0x3d,0x5c,
0x9d,0xb1,0x36,0x06,0xc0,0x34,0x73,0xf2,0x3c,0xe3,0xfe,0x95,
0x9e,0xe9,0x31,0x73,0x92,0x8b,0x0d,0xb3,0x85,0x66,0xfb,0x73,
0xe5,0x73,0xf2,0xa6,0x3d,0x90,0x27,0xe3,0xbd,0x58,0x17,0xef,
0xe6,0x6a,0xfc,0x32,0x2a,0x07,0xe6,0x9c,0x30,0xcd,0x09,0x98,
0xae,0xaa,0x43,0x9a,0x4f,0x6b,0x33,0x55,0x88,0x98,0xe9,0x50,
0x08,0x57,0xe6,0x65,0x3b,0xdb,0xc8,0xda,0x56,0x89,0x22,0x68,
0x2f,0x26,0xf6,0x9e,0xb2,0x67,0xad,0xde,0x3d,0x44,0x89,0x93,
0x7f,0x08,0xb7,0xa4,0xaf,0x9e,0x1e,0x5e,0x3d,0xbd,0x47,0x65,
0x11,0x84,0x7e,0x50,0xf1,0xe5,0xe5,0xf3,0x03,0xfe,0xdf,0xa7,
0x3a,0x89,0x25,0xe5,0xc7,0xc4,0xbe,0xe5,0x9c,0x1c,0x5a,0x06,
0xf5,0xd0,0x96,0x60,0xb0,0xeb,0x66,0xe4,0xbd,0x41,0x60,0x61,
0xa7,0x0b,0xc1,0x1d,0x81,0x5a,0x9e,0xb7,0xe7,0xc7,0x0f,0x13,
0x9d,0x99,0xe1,0xec,0x9a,0xed,0x59,0xa4,0xd2,0xd8,0xb0,0xd7,
0xae,0x0c,0xa6,0x1e,0xd8,0x89,0x8b,0x32,0x8e,0x49,0xbf,0xf4,
0x38,0x4d,0xd6,0xb9,0xdc,0x93,0x4a,0x31,0x87,0xb1,0x23,0xbe,
0x10,0x66,0x24,0x3e,0x8a,0xad,0x3c,0x08,0xee,0xa1,0x1b,0x67,
0xec,0x99,0x5c,0x21,0xb5,0x16,0x96,0x9e,0xab,0xf1,0xf3,0x8b,
0x61,0x2f,0xa1,0x0a,0xd5,0x49,0x3e,0x34,0xc7,0x82,0xf5,0x12,
0x1b,0x8b,0xc8,0xeb,0x53,0x5c,0xd2,0x89,0x59,0xca,0x52,0x52,
0x71,0x57,0xa7,0x88,0xe4,0xc1,0xa8,0xdc,0xb1,0x34,0x45,0x07,
0x3f,0x7c,0xe0,0x9d,0x5a,0x65,0xb2,0x47,0x1e,0x69,0x3b,0xa5,
0x0a,0x2b,0x77,0x35,0xd0,0x2c,0x2d,0x4a,0xc6,0x51,0x84,0x93,
0x59,0xba,0x56,0xb1,0xdb,0x39,0xd1,0x46,0x8b,0x2b,0x4a,0x70,
0x00,0xa0,0xa1,0xbd,0x74,0x72,0xa3,0xc4,0x9c,0x00,0x4e,0xd2,
0x2d,0x22,0xb3,0xc1,0x09,0xa8,0x66,0x5e,0x52,0xd8,0xb5,0x74,
0x6b,0xbf,0x92,0x03,0xa8,0xfa,0x64,0x63,0xf1,0x17,0xf5,0xb3,
0x18,0x4a,0xc6,0xd8,0x1b,0x8d,0xc5,0x6c,0x99,0x4e,0xda,0x7d,
0x29,0xe7,0x93,0x9a,0x2b,0xa6,0xd4,0xda,0xb3,0x43,0x31,0x71,
0x07,0x3e,0x50,0x61,0xa6,0xab,0x47,0x4c,0x9f,0xce,0x8e,0x4b,
0x13,0x72,0x13,0x35,0xd0,0x8c,0x99,0xfe,0xf0,0x48,0x31,0x9e,
0x19,0x0e,0x97,0xb0,0xea,0x67,0xa9,0xff,0x2d,0x9b,0x56,0xec,
0x26,0x76,0xf9,0x02,0x5a,0xa3,0x52,0x29,0xfd,0x6a,0x42,0x28,
0x86,0x83,0x40,0x9c,0x8e,0x26,0x25,0x69,0xdd,0x59,0x95,0xe0,
0x5a,0x8e,0x13,0x56,0x06,0xf7,0xb3,0xba,0x94,0x4a,0x1e,0x8b,
0x3e,0x14,0x75,0x51,0x22,0xb2,0x02,0x00,0xe3,0x22,0x3d,0x1e,
0x0e,0xba,0x4c,0xc6,0x9a,0x9b,0x13,0xee,0x98,0xef,0x93,0xfe,
0x66,0xec,0xa5,0xe7,0xfa,0x0c,0xa7,0x33,0x99,0x0e,0x1e,0x26,
0x3e,0x30,0x84,0x4a,0xbf,0x3f,0xbc,0xc4,0x6e,0x81,0xec,0x4a,
0x02,0x76,0x70,0x06,0xea,0x08,0x90,0x70,0x5d,0xa8,0x73,0xf2,
0x7b,0xd6,0x4d,0x11,0xa1,0x6b,0x4e,0x89,0xc7,0x77,0x5c,0xd9,
0x98,0x9a,0xc1,0xdc,0x02,0xb2,0xd6,0x06,0x1b,0xaa,0xe4,0xc4,
0x77,0x8f,0xbb,0xbc,0x48,0x1e,0x6b,0xc7,0x53,0x6b,0x01,0x9c,
0x70,0x06,0xac,0x88,0xf7,0x10,0xcf,0x44,0x34,0xe5,0x80,0x25,
0xe7,0xdb,0xb1,0xc6,0x3d,0x58,0x81,0xff,0x87,0x26,0xde,0xd9,
0xd4,0xb1,0x37,0x92,0x43,0xbd,0xa4,0x2f,0x44,0x73,0x97,0xde,
0x5e,0xfe,0x01,0x9a,0x19,0x88,0x3b,0x35,0x21,0xff,0x7e,0x2d,
0xb1,0x21,0xd2,0x34,0xf0,0xd2,0x05,0x53,0x72,0xa7,0xe7,0x24,
0x0c,0x8d,0xab,0x1c,0x26,0xc7,0xd9,0x05,0xd8,0xca,0xa0,0x0b,
0xd7,0x55,0x95,0x70,0xe3,0xce,0x88,0xfe,0x54,0x4a,0x56,0xdf,
0x8e,0x0e,0xa5,0x1e,0xc2,0x84,0x85,0xde,0x58,0x07,0xb6,0xe0,
0x57,0x61,0xfc,0x9d,0xfa,0x1d,0x9a,0x62,0xed,0x39,0x63,0x1c,
0xca,0xde,0x7c,0x38,0x7d,0x82,0xdf,0x18,0x98,0xcc,0xdb,0x87,
0x0d,0x41,0x0e,0x01,0x67,0x00,0x89,0xb0,0xc0,0x85,0xc7,0x28,
0xa2,0x38,0x1b,0x95,0x2d,0x1e,0x71,0x71,0x8e,0xfb,0xd9,0xd8,
0x13,0xad,0xd8,0x4b,0xd3,0xf5,0x20,0xc1,0xf8,0xe3,0xb4,0xfe,
0x1f,0x1d,0x45,0x39,0x11,0xa4,0x51,0xdc,0x4f,0xb3,0x18,0x25,
0x5d,0x3d,0x56,0xcc,0xda,0x2b,0xf2,0xcf,0xdf,0x0a,0x07,0xbb,
0x48,0x6e,0xc6,0x66,0x1f,0x19,0xd4,0xa1,0x58,0x3f,0x6a,0x4c,
0x8b,0x9b,0xdc,0xd8,0x2f,0x02,0xbb,0x48,0x90,0xc7,0xd6,0xdb,
0x1c,0xd9,0x8a,0x13,0x65,0xd4,0xbc,0xa3,0x69,0x94,0xa6,0xcc,
0xf4,0xeb,0xc9,0x8c,0x78,0x36,0x3d,0x23,0xd2,0x35,0x6d,0xc4,
0x7c,0x80,0x95,0x8b,0x13,0x35,0x26,0x7c,0x2a,0x84,0x23,0x3d,
0x13,0x5f,0x6c,0xf1,0xaa,0x2a,0x2c,0x12,0xd3,0xe2,0x22,0x1b,
0x8f,0xf1,0xf8,0x29,0x87,0x19,0x98,0x9d,0x14,0xe9,0xd2,0x7a,
0xd7,0x6d,0x22,0x88,0xe4,0x9f,0x00,0xec,0x9a,0x1d,0xd2,0xed,
0x9f,0xb5,0xd2,0x85,0xcb,0xb3,0xa7,0x61,0x22,0x37,0xae,0xbc,
0xfe,0x8e,0x39,0xbe,0x2d,0x2f,0x67,0xaf,0x66,0x97,0xe9,0x25,
0x5c,0x9a,0xc7,0xb7,0x90,0xbb,0x4d,0xa2,0x77,0x17,0x99,0x6c,
0xd9,0xa3,0x33,0x7e,0xe2,0xfd,0x1b,0x00,0x00,0xff,0xff,0x5f,
0x29,0xc9,0xd1,0xf3,0x0d,0x00,0x00,
	}))

	if err != nil {
		panic("Decompression failed: " + err.Error())
	}

	var b bytes.Buffer
	io.Copy(&b, gz)
	gz.Close()

	return b.Bytes()
}