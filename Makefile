build:
	@gomobile bind -o ./target/ClashKit.xcframework -target=ios,iossimulator,macos -ldflags=-w ./clash
