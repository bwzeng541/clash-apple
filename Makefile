build:
	@gomobile bind -o ./ClashKit.xcframework -target=ios,iossimulator,maccatalyst,macos -ldflags=-w ./
