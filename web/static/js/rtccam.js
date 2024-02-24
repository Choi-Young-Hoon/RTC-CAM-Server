
let myClientId = 0;
let rtcCamSocket = null;

document.addEventListener('DOMContentLoaded', function() {
    rtcCamSocket = new WebSocket('ws://' + window.location.hostname + ':50001/rtccam');

    rtcCamSocket.onopen = function() {
        console.log("WebSocket opened");
    }

    rtcCamSocket.onerror = function(event) {
        alert("WebSocket error");
    }

    rtcCamSocket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        console.log(data);

        if (data.client_id !== undefined) {
            handlerConnectMessage(data);
        } else if (data.is_success !== undefined) {
            handlerResultMessage(data);
        }
    }


    document.querySelector('button[data-bs-target="#roomListMenu"]').click();
});

function handlerConnectMessage(data) {
    console.log("client_id: " + data.client_id);
    myClientId = data.client_id;
    requestRoomList();
}

function handlerResultMessage(data) {
    if (data.is_success === true) {
        showRoomList(data.rooms);
    } else {

    }
}

function requestRoomList() {
    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'room_list',
        }
    }));
}

function requestRoomJoin(roomId) {
    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'join_room',
            join_room_id: roomId,
        }
    }));
}

function requestRoomLeave() {
    rtcCamSocket.send(JSON.stringify({
        type: 'leave_room',
    }));
}

function showRoomList(roomList) {
    let roomListElement = document.getElementById('roomList');
    roomListElement.innerHTML = '';

    for (let roomId in roomList.rooms) {
        let room = roomList.rooms[roomId];

        let clientCount = Object.keys(room.clients).length === 0 ? 0 : Object.keys(room.clients).length;
        let htmlString = "<h3>" + room.title + "   이용자 수: " + clientCount + "</h3><hr>";

        let roomElement = document.createElement('div');
        roomElement.onmouseover = function() {
            roomElement.style.backgroundColor = "lightgray";
            roomElement.style.color = "blue";
        }
        roomElement.onmouseleave = function() {
            roomElement.style.backgroundColor = "white";
            roomElement.style.color = "black";
        }
        roomElement.onmousedown = function() {
            roomElement.style.backgroundColor = "gray";
        }
        roomElement.onclick = function() {
            requestRoomJoin(room.id);
        }
        roomElement.style.cursor = "pointer";
        roomElement.innerHTML = htmlString;
        roomListElement.appendChild(roomElement);
    }
}