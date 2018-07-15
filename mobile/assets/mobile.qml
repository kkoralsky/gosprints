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
    title: "Gosprints control app"

    property int tournamentPlayers: 4
    property alias tabBar: tabBar

    Component.onCompleted: {
        Qt.application.name = "app name"
        Qt.application.domain = "domain"
        Qt.application.organization = "organization"

        settingsPopup.open()
    }

    Settings {
        id: settings
        property alias connectionHost: connectionHostField.text
        property alias connectionPort: connectionPortField.text
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
                    width: settingsPopup.width / 3
                }

                TextField {
                    id: connectionHostField
                    placeholderText: connectionHost
                    width: settingsPopup.width / 2
                }
            }

            RowLayout {
                spacing: 10

                Label {
                    text: "Port:"
                    width: settingsPopup.width / 3
                }

                TextField {
                    id: connectionPortField
                    placeholderText: connectionPort
                    width: settingsPopup.width / 2
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
                                true  // block
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
        id: aboutDialog
        modal: true
        focus: true
        x: (window.width - width) / 2
        y: window.height / 6
        width: Math.min(window.width, window.height) / 3 * 2
        contentHeight: aboutColumn.height

        Column {
            id: aboutColumn
            spacing: 20

            Label {
                text: "About"
                font.bold: true
            }

            Label {
                width: aboutDialog.availableWidth
                text: "The Qt Quick Controls 2 module delivers the next generation user interface controls based on Qt Quick."
                wrapMode: Label.Wrap
                font.pixelSize: 12
            }

            Label {
                width: aboutDialog.availableWidth
                text: "In comparison to the desktop-oriented Qt Quick Controls 1, Qt Quick Controls 2 "
                    + "are an order of magnitude simpler, lighter and faster, and are primarily targeted "
                    + "towards embedded and mobile platforms."
                wrapMode: Label.Wrap
                font.pixelSize: 12
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
        anchors.fill: parent
        color: "black"
        opacity: 0.3
        BusyIndicator {
            running: ! (SprintsClient.connState === 0 || SprintsClient.connState === 2)
            anchors.centerIn: parent
            onRunningChanged: console.log(SprintsClient.connState)
        }
    }
}
