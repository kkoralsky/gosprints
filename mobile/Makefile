
ANDROID_ARTIFACT=./deploy/android/build-debug.apk
LINUX_ARTIFACT=./deploy/linux/mobile

.PHONY: qmlscene runqmldebug runqml moc\
	deploylinux deployandroid runandroid

mobile: main.go client.go results.go
	go build

moc: moc.go

deploylinux: $(LINUX_ARTIFACT)

deployandroid: $(ANDROID_ARTIFACT)

runandroid: .runandroid
	adb install -sr $(ANDROID_ARTIFACT)

qmlscene: 
	$(QT_DIR)/$(QT_VERSION)/gcc_64/bin/qmlscene assets/App.qml

runqmldebug: .qmldebug
	./mobile -qml_debug 

runqml: .qmldebug rcc.cpp
	./mobile

moc.go: client.go results.go client_mock.go
	qtmoc -debug

rcc.cpp: rcc.qrc assets/**
	qtrcc -debug

$(LINUX_ARTIFACT): *.go assets/**
	qtdeploy -debug build linux
	@touch mobile  # ^ modifies moc.go

$(ANDROID_ARTIFACT): *.go rcc.qrc assets/** android/**
	qtdeploy -debug build android
	@touch mobile  # ^ modifies moc.go

.qmldebug: mobile
	@touch $@

.runandroid: $(ANDROID_ARTIFACT)
	@touch $@

updateqt:
	go get -v github.com/therecipe/qt/cmd/... 
	CGO_CXXFLAGS_ALLOW=".*" CGO_LDFLAGS_ALLOW=".*" CGO_CFLAGS_ALLOW=".*" $$GOPATH/bin/qtsetup

include Makefile.custom
