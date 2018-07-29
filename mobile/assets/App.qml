import QtQuick 2.6
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.2
import QtQuick.Controls.Material 2.0
import QtQuick.Controls.Universal 2.0
import Qt.labs.settings 1.0
import "FontAwesome"


ApplicationWindow {
    id: window
    width: 360
    height: 520
    visible: true
    // title: "Gosprints control app"

    property int tournamentPlayers: 4
    property alias tabBar: tabBar

    enum ConnState {
        Idle,
        Connecting,
        Ready,
        TransientFailure,
        Shutdown
    }

    Component.onCompleted: {
        Qt.application.name = "app name"
        Qt.application.domain = "domain"
        Qt.application.organization = "organization"

        SprintsClient.dialGrpc(settings.connectionHost, settings.connectionPort, false)
    }

    Settings {
        id: settings
        property alias connectionHost: connectionHostField.text
        property alias connectionPort: connectionPortField.text
    }

    Connections {
        target: SprintsClient

        onConnStateChanged: {
            if(SprintsClient.connState === App.ConnState.Ready) {
                connectingOverlay.visible=false
                settings.connectionHost=connectionHostField.text
                settings.connectionPort=connectionPortField.text
            } else {
                connectingOverlay.visible=true
            } if (SprintsClient.connState === App.ConnState.TransientFailure) {
                settingsPopup.open()
            }
        }
    }

    header: ToolBar {

        Material.foreground: "white"

        RowLayout {
            spacing: 20
            anchors.fill: parent

            ToolButton {
                contentItem: Text {
                    font.family: FontAwesome.fontFamily
                    text: FontAwesome.plug
                    color: "#ccc"
                    font.pixelSize: 30
                }

                onClicked: {
                    settingsPopup.open()
                }
            }

            Label {
                id: titleLabel
                text: "Gosprints Ctrl"
                font.pixelSize: 20
                elide: Label.ElideRight
                horizontalAlignment: Qt.AlignHCenter
                verticalAlignment: Qt.AlignVCenter
                Layout.fillWidth: true
            }

            ToolButton {
                contentItem: Text {
                    font.family: FontAwesome.fontFamily
                    text: FontAwesome.ellipsisV
                    font.pixelSize: 30
                }
                onClicked: optionsMenu.open()

                Menu {
                    id: optionsMenu
                    x: parent.width - width
                    transformOrigin: Menu.TopRight

                    MenuItem {
                        text: "New Tournament"
                        onTriggered: {
                            stackView.clear()
                            stackView.push(newTournamentPage)
                        }
                    }

                    MenuItem {
                        text: "Load Tournament"
                        onTriggered: {
                            loadTournamentPopup.open()
                        }
                    }
                }
            }
        }
    }

    StackView {
        id: stackView
        anchors.fill: parent
        anchors.margins: 20

        initialItem: newTournamentPage

        Tournament {
            id: tournamentPage
        }

        NewTournament {
            id: newTournamentPage
        }
    }

    Popup {
        id: settingsPopup
        x: (window.width - width) / 2
        y: window.height / 6
        width: Math.min(window.width, window.height) / 3 * 2
        height: settingsColumn.implicitHeight + topPadding + bottomPadding
        modal: true
        focus: true

        contentItem: ColumnLayout {
            id: settingsColumn
            spacing: 20

            Label {
                text: "Connection Settings"
                font.bold: true
            }

            RowLayout {
                spacing: 10

                Label {
                    text: "Host:"
                    // width: settingsPopup.width / 3
                }

                TextField {
                    id: connectionHostField
                    placeholderText: connectionHost
                    // width: settingsPopup.width / 2
                }
            }

            RowLayout {
                spacing: 10

                Label {
                    text: "Port:"
                    // width: settingsPopup.width / 3
                }

                TextField {
                    id: connectionPortField
                    placeholderText: connectionPort
                    // width: settingsPopup.width / 2
                }
            }

            RowLayout {
                spacing: 10

                Button {
                    id: connectButton
                    text: "Connect"
                    onClicked: {
                        SprintsClient.dialGrpc(
                                connectionHostField.text,
                                parseInt(connectionPortField.text),
                                false  // dont block
                        )
                        settingsPopup.close()
                    }

                    Material.foreground: Material.primary
                    Material.background: "transparent"
                    Material.elevation: 0

                    Layout.preferredWidth: 0
                    Layout.fillWidth: true
                }

                Button {
                    id: cancelButton
                    text: "Cancel"
                    onClicked: {
                        settingsPopup.close()
                    }

                    Material.background: "transparent"
                    Material.elevation: 0

                    Layout.preferredWidth: 0
                    Layout.fillWidth: true
                }
            }
        }
    }

    Popup {
        id: loadTournamentPopup
        modal: true
        focus: true
        x: (window.width - width) / 2
        y: window.height / 6
        width: Math.min(window.width, window.height) / 3 * 2
        contentHeight: loadTournamentColumn.height

        Connections {
            target: TournamentConfig
            onCurrentIndexChanged: {
                // console.log(TournamentConfig.currentIndex)
                loadTournamentCombo.currentIndex = TournamentConfig.currentIndex
            } 
        }

        Column {
            id: loadTournamentColumn
            width: parent.width
            spacing: 20

            Label {
                text: "Load Tournament"
                font.bold: true
            }

            ComboBox {
                id: loadTournamentCombo
                model: TournamentConfig.tournaments
                width: parent.width
                currentIndex: TournamentConfig.currentIndex
                onActivated: {
                    loadTournamentPopup.close()
                    SprintsClient.loadTournament(model[index])
                    SprintsClient.getResults("MALE")
                }
            }
        }
    }

    footer: TabBar {
        id: tabBar
        currentIndex: tournamentPage.currentIndex

        TabButton {
            text: "Race"
            onClicked: {
                stackView.clear()
                stackView.push(tournamentPage)
            }
        }
        TabButton {
            text: "Results"
            onClicked: {
                stackView.clear()
                stackView.push(tournamentPage)
                SprintsClient.getResults("MALE")
            }
        }
    }

    Rectangle {
        id: connectingOverlay
        anchors.fill: parent
        color: "black"
        opacity: 0.3
        BusyIndicator {
            anchors.centerIn: parent
        }
    }
}
