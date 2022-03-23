build:
	@gomobile bind -o ./ClashKit.xcframework -target=ios,macos -ldflags=-w ./
