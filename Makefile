build:
	@gomobile bind -o ./target/ClashKit.xcframework -target=macos -ldflags=-w ./clash
