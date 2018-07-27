package main

import (
	"flag"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	gui.NewQGuiApplication(0, []string{})

	var (
		connectionHost = core.NewQVariant14("localhost")
		connectionPort = core.NewQVariant8(8888)

		resultModel   = NewResultModel(nil)
		engine        = qml.NewQQmlApplicationEngine(nil)
		root          = engine.RootContext()
		sprintsClient = setupMockSprintsClient(resultModel)
	)

	root.SetContextProperty("SprintsClient", sprintsClient)
	root.SetContextProperty("ResultModel", sprintsClient.resultModel)
	root.SetContextProperty2("connectionHost", connectionHost)
	root.SetContextProperty2("connectionPort", connectionPort)

	engine.Load(core.NewQUrl3("./assets/App.qml", 0))

	gui.QGuiApplication_Exec()

	m.Run()
}
