import QtQuick 2.6
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.2

Pane {
    ListView {
        id: listView
        anchors.fill: parent
        model: ResultModel
        Layout.fillWidth: true

        Layout.fillHeight: true
        clip: true

        section.property: "destValue" 
        section.delegate: Pane {

            width: listView.width
            height: sectionLabel.implicitHeight + 20

            Label {
                id: sectionLabel
                text: section + (newTournamentPage.mode == NewTournament.TournamentMode.DISTANCE ? "m" : "s")
                anchors.centerIn: parent
            }
        }

        delegate: ItemDelegate {
            width: parent.width

            RowLayout {
                anchors.fill: parent
                Label  {
                    text: name
                }

                Label {
                    Layout.fillWidth: true
                    text: {
                        if(newTournamentPage.mode == NewTournament.TournamentMode.DISTANCE) {
                            return score / 1000 + "s"
                        } else {
                            return score + "m"
                        }
                    }
                    horizontalAlignment: Qt.AlignRight
                }
            }
        }
    }
}
