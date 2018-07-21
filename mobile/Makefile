
ANDROID_ARTIFACT=./deploy/android/build-debug.apk
LINUX_ARTIFACT=./deploy/linux/mobile

.PHONY: qmlscene runqmldebug runqml\
	deploylinux deployandroid runandroid

mobile: main.go moc.go
	go build

deploylinux: $(LINUX_ARTIFACT)

deployandroid: $(ANDROID_ARTIFACT)

runandroid: deployandroid
	adb install -sr $(ANDROID_ARTIFACT)

qmlscene: 
	$(QT_DIR)/$(QT_VERSION)/gcc_64/bin/qmlscene assets/mobile.qml

runqmldebug: .qmldebug
	./mobile -qml_debug 

runqml: .qmldebug rcc.cpp
	./mobile

moc.go: client.go results.go
	qtmoc -debug

rcc.cpp: rcc.qrc assets/**
	qtrcc -debug

$(LINUX_ARTIFACT): mobile
	qtdeploy -debug build linux

$(ANDROID_ARTIFACT): mobile
	qtdeploy -debug build android

.qmldebug: mobile
	@touch .qmldebug
