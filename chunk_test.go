package piio

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var uncompressedPi = []byte{
	3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5, 9,
}

var compressedPi = []byte{
	0x31, 0x41, 0x59, 0x26, 0x53, 0x59,
}

var textPi = "314159265359"

func TestCompressedDigits(t *testing.T) {
	var chnk *CompressedChunk
	Convey("Given a CompressedChunk", t, func() {
		chnk = &CompressedChunk{
			firstIndex: 0,
			data:       compressedPi,
		}

		Convey("Digits() should", func() {
			Convey("error on out of range inputs.", func() {
				var err error
				_, err = chnk.Digit(-1)
				So(err, ShouldNotBeNil)
				_, err = chnk.Digit(12)
				So(err, ShouldNotBeNil)
			})
			Convey("work on in range inputs.", func() {
				var err error
				var b byte

				b, err = chnk.Digit(0)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 3)

				b, err = chnk.Digit(7)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 6)

				b, err = chnk.Digit(11)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 9)
			})
		})
	})
	Convey("Given a CompressedChunk", t, func() {
		chnk = &CompressedChunk{
			firstIndex: 2,
			data:       compressedPi[1:],
		}

		Convey("Digits() should", func() {
			Convey("error on out of range inputs.", func() {
				var err error
				_, err = chnk.Digit(1)
				So(err, ShouldNotBeNil)
				_, err = chnk.Digit(12)
				So(err, ShouldNotBeNil)
			})
			Convey("work on in range inputs.", func() {
				var err error
				var b byte

				b, err = chnk.Digit(2)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 4)

				b, err = chnk.Digit(7)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 6)

				b, err = chnk.Digit(11)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 9)
			})
		})
	})
}

func TestUncompressedDigits(t *testing.T) {
	var chnk *UncompressedChunk
	Convey("Given an UncompressedChunk", t, func() {
		chnk = &UncompressedChunk{
			FirstDigitIndex: 0,
			Digits:          uncompressedPi,
		}

		Convey("Digits() should", func() {
			Convey("error on out of range inputs.", func() {
				var err error
				_, err = chnk.Digit(-1)
				So(err, ShouldNotBeNil)
				_, err = chnk.Digit(12)
				So(err, ShouldNotBeNil)
			})
			Convey("work on in range inputs.", func() {
				var err error
				var b byte

				b, err = chnk.Digit(0)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 3)

				b, err = chnk.Digit(7)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 6)

				b, err = chnk.Digit(11)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 9)
			})
		})
	})
	Convey("Given an UncompressedChunk", t, func() {
		chnk = &UncompressedChunk{
			FirstDigitIndex: 2,
			Digits:          uncompressedPi[2:],
		}

		Convey("Digits() should", func() {
			Convey("error on out of range inputs.", func() {
				var err error
				_, err = chnk.Digit(1)
				So(err, ShouldNotBeNil)
				_, err = chnk.Digit(12)
				So(err, ShouldNotBeNil)
			})
			Convey("work on in range inputs.", func() {
				var err error
				var b byte

				b, err = chnk.Digit(2)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 4)

				b, err = chnk.Digit(7)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 6)

				b, err = chnk.Digit(11)
				So(err, ShouldBeNil)
				So(b, ShouldEqual, 9)
			})
		})
	})
}

func TestCompressedBasicGetters(t *testing.T) {
	Convey("Given a basic CompressedChunk", t, func() {
		chnk := &CompressedChunk{
			firstIndex: 42,
			data:       make([]byte, 12),
		}
		Convey("the is compressed getter should work.", func() {
			So(chnk.IsCompressed(), ShouldBeTrue)
		})
		Convey("the first index getter should work.", func() {
			So(chnk.FirstIndex(), ShouldEqual, 42)
		})
		Convey("the length getter should work.", func() {
			So(chnk.Length(), ShouldEqual, 24)
		})
		Convey("the last index getter should work.", func() {
			So(chnk.LastIndex(), ShouldEqual, 42+24-1)
		})
	})
}

func TestUncompressedBasicGetters(t *testing.T) {
	Convey("Given a basic UncompressedChunk", t, func() {
		chnk := &UncompressedChunk{
			FirstDigitIndex: 42,
			Digits:          make([]byte, 24),
		}
		Convey("the is compressed getter should work.", func() {
			So(chnk.IsCompressed(), ShouldBeFalse)
		})
		Convey("the first index getter should work.", func() {
			So(chnk.FirstIndex(), ShouldEqual, 42)
		})
		Convey("the length getter should work.", func() {
			So(chnk.Length(), ShouldEqual, 24)
		})
		Convey("the last index getter should work.", func() {
			So(chnk.LastIndex(), ShouldEqual, 42+24-1)
		})
	})
}

func TestCompress(t *testing.T) {
	Convey("Given an UncompressedChunk", t, func() {
		chnk := &UncompressedChunk{
			FirstDigitIndex: 0,
			Digits:          uncompressedPi,
		}
		Convey("compressing should work.", func() {
			c := Compress(chnk)
			cc, ok := c.(*CompressedChunk)
			So(ok, ShouldBeTrue)
			So(cc.firstIndex, ShouldEqual, 0)
			So(cc.data, ShouldResemble, compressedPi)
		})
	})
	Convey("Given an CompressedChunk", t, func() {
		chnk := &CompressedChunk{
			firstIndex: 0,
			data:       compressedPi,
		}
		Convey("compressing shouldn't do anything.", func() {
			c := Compress(chnk)
			So(c, ShouldEqual, chnk)
		})
	})
}

func TestUncompress(t *testing.T) {
	Convey("Given an CompressedChunk", t, func() {
		chnk := &CompressedChunk{
			firstIndex: 0,
			data:       compressedPi,
		}
		Convey("decompressing should work.", func() {
			c := Decompress(chnk)
			cc, ok := c.(*UncompressedChunk)
			So(ok, ShouldBeTrue)
			So(cc.FirstDigitIndex, ShouldEqual, 0)
			So(cc.Digits, ShouldResemble, uncompressedPi)
		})
	})
	Convey("Given an UncompressedChunk", t, func() {
		chnk := &UncompressedChunk{
			FirstDigitIndex: 0,
			Digits:          uncompressedPi,
		}
		Convey("decompressing shouldn't do anything.", func() {
			c := Decompress(chnk)
			So(c, ShouldEqual, chnk)
		})
	})
}
