import QtQuick 2.6
import QtQuick.Layouts 1.3
import QtQuick.Controls 2.2


Page {
    id: tournamentPage

    property alias currentIndex: tournamentSwipeView.currentIndex

    SwipeView {
        anchors.fill: parent
        id: tournamentSwipeView
        currentIndex: tabBar.currentIndex

        Race {
            id: racePane
        }

        Pane {
            id: resultsPane
            // anchors.fill: parent

            ListView {
                anchors.fill: parent
                model: ResultModel
                delegate: Text {
                    text: "score:" + score + "  name:" + name
                }
            }
        }
    }
}
