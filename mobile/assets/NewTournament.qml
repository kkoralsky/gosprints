import QtQuick 2.6
import QtQuick.Controls 2.0
import QtQuick.Dialogs 1.0

Page {
    id: newTournamentPage

    property string name: tournamentNameTextField.text
    property int playerCount: tournamentPlayerCountSpinBox.value
    property string mode: "D" // or "T" for time
    property int destValue: tournamentDesintationSpinBox.value

    Binding {
       target: newTournamentPage
       property: "playerCount"
       value: parseInt(tournamentPlayerCountTextField.text)
    }

    Binding {
        target: newTournamentPage
        property: "mode"
        value: tournamentDistanceRadio.checked ? "D" : "T"
    }

    Binding {
        target: newTournamentPage
        property: "destValue"
        value: parseInt(tournament)
    }

    Column {
        spacing: 10
        Row {
            Label {
                text: "Tournament name"
                width: newTournamentPage.width / 3
                height: 40
            }
            TextField {
                id: tournamentNameTextField
                width: newTournamentPage.width / 2
                height: 40
                placeholderText: "Goldpsrints"
            }
        }
        Row {
            Label {
                text: "Player count"
                height: 40
                width: newTournamentPage.width / 3
            }
            SpinBox {
                id: tournamentPlayerCountSpinBox
                width: newTournamentPage.width / 2
                height: 40
                value: 2
                from: 1
                to: 10
            }
        }
        Row {
            Label {
                text: "Race mode"
                verticalAlignment: "AlignVCenter"
                width: newTournamentPage.width / 3
            }
            RadioButton {
                id: tournamentDistanceRadio
                text: "distance"
                checked: true
            }
            RadioButton {
                id: tournamentTimeRadio
                text: "time"
                checked: false
            }
        }
        Row {
            spacing: 5
            Label {
                text: tournamentDistanceRadio.checked ? "Distance" : "Time"
                verticalAlignment: "AlignVCenter"
                height: 40
                width: newTournamentPage.width / 3
            }
            SpinBox {
                id: tournamentDesintationSpinBox
                width: newTournamentPage.width / 2
                height: 40
                value: tournamentDistanceRadio.checked ? 400 : 25
                from: tournamentDistanceRadio.checked ? 50 : 10
                to: tournamentDistanceRadio.checked ? 4000 : 5*60
            }
            Label {
                text: tournamentDistanceRadio.checked ? "m" : "s"
                verticalAlignment: "AlignVCenter"
                height: 40
            }
        }

        Button {
            text: "Setup"
            onClicked: SprintsClient.newTournament(
                newTournamentPage.name,
                newTournamentPage.destValue,
                newTournamentPage.mode,
                newTournamentPage.playerCount,
                ["blue", "red", "green", "yellow", "white", "rose", "brown",
                 "orange", "gray"]  // hardcoded color names FIXME 
            )
        }
    }
}
