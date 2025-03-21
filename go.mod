module github.com/vbsw/g2d

go 1.13

require (
	github.com/vbsw/golib/cdata v0.4.0
	github.com/vbsw/g2d/dummycontext v0.1.0
	github.com/vbsw/g2d/window v0.1.0
)

replace (
	github.com/vbsw/golib/cdata => ../golib/cdata
	github.com/vbsw/g2d/dummycontext => ./dummycontext
	github.com/vbsw/g2d/window => ./window
)
