package main

import (
	"flag"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
)

var (
	defaultHost      = "localhost"
	defaultPort uint = 9998
)

func main() {
	var (
		qmlPath = "assets/App.qml"
	)
	qmlDebug := flag.Bool("qml_debug", false, "load qml from filesystem rather than from QRC")
	flag.Parse()

	if !*qmlDebug {
		qmlPath = "qrc:/" + qmlPath
	}

	core.QCoreApplication_SetApplicationName("Gosprints Ctrl")
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	gui.NewQGuiApplication(len(os.Args), os.Args)

	var (
		connectionHost = core.NewQVariant14(defaultHost)
		connectionPort = core.NewQVariant8(defaultPort)

		resultModel   = NewResultModel(nil)
		engine        = qml.NewQQmlApplicationEngine(nil)
		root          = engine.RootContext()
		sprintsClient = SetupSprintsClient(resultModel)
	)

	root.SetContextProperty("SprintsClient", sprintsClient)
	root.SetContextProperty2("connectionHost", connectionHost)
	root.SetContextProperty2("connectionPort", connectionPort)

	engine.Load(core.NewQUrl3(qmlPath, 0))

	gui.QGuiApplication_Exec()

	sprintsClient.Close()
}
