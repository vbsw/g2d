module github.com/vbsw/g2d

go 1.13

require (
	github.com/vbsw/oglr v0.0.0-20220616120830-2cb9df765214
	github.com/vbsw/oglwnd v0.1.2
)

replace (
	github.com/vbsw/oglr => ../oglr
	github.com/vbsw/oglwnd => ../oglwnd
)
