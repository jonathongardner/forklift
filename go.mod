module github.com/jonathongardner/forklift

go 1.23

toolchain go1.23.3

require (
	github.com/diskfs/go-diskfs v1.5.0
	github.com/gabriel-vasile/mimetype v1.4.8
	github.com/jonathongardner/libarchive v0.0.0-20240814185856-047f3efea4c8
	github.com/jonathongardner/virtualfs v0.0.2
	github.com/sirupsen/logrus v1.9.4-0.20230606125235-dd1b4c2e81af
	github.com/sorairolake/lzip-go v0.3.5
	github.com/ulikunitz/xz v0.5.12
	github.com/urfave/cli/v2 v2.27.3
	kraftkit.sh v0.9.4
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/djherbis/times v1.6.0 // indirect
	github.com/elliotwutingfeng/asciiset v0.0.0-20230602022725-51bbb787efab // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/pierrec/lz4/v4 v4.1.21 // indirect
	github.com/pkg/xattr v0.4.9 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/net v0.34.0 // indirect
	golang.org/x/sys v0.29.0 // indirect
)

replace github.com/diskfs/go-diskfs => github.com/jonathongardner/go-diskfs v0.0.0-20250125020509-6251481bca9f
