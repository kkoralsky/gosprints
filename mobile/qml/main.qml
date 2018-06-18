import QtQuick 2.0

Rectangle
{
	width: 320
	height: 240

	Text
	{
        id: textView
		anchors.centerIn: parent
		text: "click here!"

		Connections
		{
			target: SprintsClient

			onInfo: textView.text = msg
			onError: textView.text = "Failed to " + err + " for " + addr + "\nWith error message: " + msg + "\n"
			onSuccess: textView.text = "Received: \"" + msg + "\"\n"
		}
    }

	MouseArea
	{
		anchors.fill: parent
		onClicked: SprintsClient.newRace(["hello", "world"], 30)
	}
    Component.onCompleted: console.log(SprintsClient.dialGrpc("localhost", 9999))
}
