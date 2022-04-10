build:
	@gomobile bind -o ./ClashKit.xcframework -target=ios,iossimulator,macos -ldflags=-w ./clash
