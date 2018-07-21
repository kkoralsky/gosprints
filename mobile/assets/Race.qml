import QtQuick 2.6
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.2


Pane {
    enum State {
        Preparing,
        Starting,
        Racing
    }
    id: pane

    property int raceState: Race.State.Preparing;

    Column {
        spacing: 10

        Repeater {
            id: newRaceRepeater
            model: newTournamentPage.playerCount
            Row {
                Label {
                    width: pane.width / 3
                    height: 40
                    text: "player #"+(index+1)
                }

                TextField {
                    width: pane.width / 2
                    height: 40
                    placeholderText: "player name"
                }
            }
        }

        Row {
            Label {
                width: pane.width / 3
                height: 40
                verticalAlignment: "AlignVCenter"
                text: newTournamentPage.mode == "D" ? "Distance" : "Time"
            }
            SpinBox {
                width: pane.width / 2
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

        RowLayout {
            spacing: 20
            anchors.horizontalCenter: parent.horizontalCenter

            Button {
                text: pane.raceState == Race.State.Preparing ? "New race" : (pane.raceState == Race.State.Starting) ? "Start" : "Stop"
                onClicked: {
                    switch (pane.raceState) {
                        case Race.State.Preparing:
                            var racers = []
                            for(var i=0; i<newRaceRepeater.count; i++) {
                                racers.push(newRaceRepeater.itemAt(i).children[1].text);
                            }
                            pane.raceState = Race.State.Starting
                            SprintsClient.newRace(racers, 5)
                            break
                        case Race.State.Starting:
                            pane.raceState = Race.State.Racing
                            SprintsClient.startRace() 
                            break
                        case Race.State.Racing:
                            pane.raceState = Race.State.Preparing
                            SprintsClient.abortRace()
                            break
                        default:
                            break   // should never happen
                    }
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
