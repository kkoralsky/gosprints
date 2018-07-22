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
                text: newTournamentPage.mode == NewTournament.TournamentMode.DISTANCE ? "Distance" : "Time"
            }
            SpinBox {
                id: destValueSpinBox 
                width: pane.width / 2
                height: 40
                value: newTournamentPage.destValue
                from: newTournamentPage.mode == NewTournament.TournamentMode.DISTANCE ? 50 : 10
                to: newTournamentPage.mode == NewTournament.TournamentMode.DISTANCE ? 4000 : 5*60
            }
            Label {
                height: 40
                verticalAlignment: "AlignVCenter"
                text: newTournamentPage.mode == NewTournament.TournamentMode.DISTANCE ? "m" : "s"
            }
        }

        RowLayout {
            spacing: 20
            anchors.horizontalCenter: parent.horizontalCenter

            Button {
                text: "New race"
                enabled: pane.raceState != Race.State.Racing
                onClicked: {
                    var racers = []
                    for(var i=0; i<newRaceRepeater.count; i++) {
                        racers.push(newRaceRepeater.itemAt(i).children[1].text);
                    }
                    pane.raceState = Race.State.Starting
                    SprintsClient.newRace(racers, destValueSpinBox.value)
                }
            }
            Button {
                text: pane.raceState == Race.State.Racing ? "Stop" : "Start"
                enabled: pane.raceState != Race.State.Preparing
                onClicked: {
                    if(pane.raceState==Race.State.Starting) {
                        pane.raceState = Race.State.Racing
                        SprintsClient.startRace()
                    } else if (pane.raceState==Race.State.Racing) {
                        pane.raceState = Race.State.Preparing 
                        SprintsClient.abortRace()
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
