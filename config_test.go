package config

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	TestIni     = "testdata/test.ini"
	NotFoundIni = "404.ini"
	MoreIni     = "testdata/more.ini"

	DevelopEnv = "develop"
	TestingEnv = "testing"
	ProductEnv = "product"
)

func TestLoad(t *testing.T) {
	Convey("Test load", t, func() {
		Convey("ini file not found", func() {
			_, err := Load(NotFoundIni)
			So(err, ShouldNotBeNil)
		})
		Convey("ini file found", func() {
			_, err := Load(TestIni)
			So(err, ShouldBeNil)
		})
	})
}

func TestConfig_Read(t *testing.T) {
	Convey("Test read", t, func() {
		conf, _ := Load(TestIni, MoreIni)

		Convey("comment begins with # or ;", func() {
			v1 := conf.Read(DevelopEnv, "key2")
			So(v1, ShouldEqual, "")
			v2 := conf.Read(DevelopEnv, "key3")
			So(v2, ShouldEqual, "")
		})
		Convey("a line doesnot contain =", func() {
			v := conf.Read(DevelopEnv, "key4")
			So(v, ShouldEqual, "")
		})
		Convey("a key's value is empty", func() {
			v := conf.Read(DevelopEnv, "key5")
			So(v, ShouldEqual, "")
		})
		Convey(`comment with substr '\t#'`, func() {
			v := conf.Read(DevelopEnv, "key6")
			So(v, ShouldEqual, "key6")
		})
		Convey(`comment with substr ' #'`, func() {
			v := conf.Read(DevelopEnv, "key7")
			So(v, ShouldEqual, "key7")
		})
		Convey(`comment with substr ' ;'`, func() {
			v := conf.Read(DevelopEnv, "key8")
			So(v, ShouldEqual, "key8")
		})
		Convey(`comment with substr '\t//'`, func() {
			v := conf.Read(DevelopEnv, "key9")
			So(v, ShouldEqual, "key9")
		})
		Convey(`comment with substr ' //'`, func() {
			v := conf.Read(DevelopEnv, "key10")
			So(v, ShouldEqual, "key10")
		})
		Convey("normal ", func() {
			v := conf.Read(DevelopEnv, "key11")
			So(v, ShouldEqual, "key11#$*()@###")
		})
		Convey("different value with same key in different sections", func() {
			v1 := conf.Read(DevelopEnv, "key1")
			So(v1, ShouldEqual, "develop")
			v2 := conf.Read(TestingEnv, "key1")
			So(v2, ShouldEqual, "testing")
			v3 := conf.Read(ProductEnv, "key1")
			So(v3, ShouldEqual, "product")
		})
		Convey("values in multiple files", func() {
			v1 := conf.Read(DevelopEnv, "key1")
			So(v1, ShouldEqual, "develop")
			v2 := conf.Read(DevelopEnv, "name")
			So(v2, ShouldEqual, "star")
			v3 := conf.Read(TestingEnv, "name")
			So(v3, ShouldEqual, "kimi")
			v4 := conf.Read(ProductEnv, "name")
			So(v4, ShouldEqual, "Kimi.Wang")
		})
	})
}
