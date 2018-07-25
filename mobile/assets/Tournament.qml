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

        Results {
            id: resultsPane

        }
    }
}
