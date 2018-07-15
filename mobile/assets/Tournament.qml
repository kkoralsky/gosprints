import QtQuick 2.6
import QtQuick.Controls 2.2

Page {
    id: tournamentPage

    property alias currentIndex: tournamentSwipeView.currentIndex

    SwipeView {
        anchors.fill: parent
        id: tournamentSwipeView
        currentIndex: tabBar.currentIndex

        Pane {
            id: newRacePane
            Column {
                spacing: 10

                Repeater {
                    id: newRaceRepeater
                    model: newTournamentPage.playerCount
                    Row {
                        Label {
                            width: newRacePane.width / 3
                            height: 40
                            text: "player #"+(index+1)
                        }

                        TextField {
                            width: newRacePane.width / 2
                            height: 40
                            placeholderText: "player name"
                        }
                    }
                }

                Row {
                    Label {
                        width: newRacePane.width / 3
                        height: 40
                        verticalAlignment: "AlignVCenter"
                        text: newTournamentPage.mode == "D" ? "Distance" : "Time"
                    }
                    SpinBox {
                        width: newRacePane.width / 2
                        height: 40
                        value: newTournamentPage.destValue
                        from: newTournamentPage.mode == "D"  ? 50 : 10
                        to: newTournamentPage.mode == "D" ? 4000 : 5*60
                    }
                    Label {
                        height: 40
                        verticalAlignment: "AlignVCenter"
                        text: newTournamentPage.mode == "D" ? "m" : "s"
                    }
                }

                Row {
                    spacing: 20

                    Button {
                        text: "New race"
                        onClicked: {
                            var racers = []
                            for(var i=0; i<newRaceRepeater.count; i++) {
                                racers.push(newRaceRepeater.itemAt(i).children[1].text);
                            }
                            SprintsClient.newRace(racers, 5)
                        }
                    }
                    Button {
                        text: "Swap"
                        onClicked: {
                            var toSwap = ""
                            for(var i=0; i<newRaceRepeater.count; i++) {
                                var textField = newRaceRepeater.itemAt(i).children[1],
                                    previous = textField.text
                                textField.text = toSwap
                                toSwap = previous
                            }
                            newRaceRepeater.itemAt(0).children[1].text = toSwap
                        }
                    }
                    Button {
                        text: "Clear"
                        onClicked: {
                            for(var i=0; i<newRaceRepeater.count; i++) {
                                var textField = newRaceRepeater.itemAt(i).children[1]
                                textField.text = ""
                            }
                        }
                    }
                }
            }
        }

        Pane {
            id: resultsPane
            Rectangle {
                color: "blue"
                anchors.fill: parent
            }
        }

    }
}
