let rtcCamSocket = null;
let rtcCamWSServerUrl = 'ws://' + window.location.hostname + ':40001/rtccam';

let myClientId = 0;
let stunServerUrl = null;
let turnServerUrl = null;

let roomClients = [];
let localPeerConnection;
let peerConnectionMap = new Map();
let peerVideoStreamMap = new Map();

let localVideoElement = document.getElementById('localVideo');
let localMediaStream = null;

document.addEventListener('DOMContentLoaded', function() {
    rtcCamSocket = new WebSocket(rtcCamWSServerUrl);

    rtcCamSocket.onopen = function() {
        console.log("WebSocket opened");
    }

    rtcCamSocket.onerror = function(event) {
        alert("rtccam 서버와 통신할 수 없습니다.");
        rtcCamSocket = new WebSocket(rtcCamWSServerUrl);
    }

    rtcCamSocket.onmessage = function(event) {
        const data = JSON.parse(event.data);
        console.log(data);

        if (data.client_id !== undefined) {
            handlerConnectMessage(data);
        } else if (data.result_message !== undefined) {
            handlerResultMessage(data);
        } else if (data.request_type !== undefined) {
            if (data.request_type === 'offer') {
                handlerOfferMessage(data);
            } else if (data.request_type === 'answer') {
                handlerAnswerMessage(data);
            } else if (data.request_type === 'candidate') {
                handlerCandidateMessage(data);
            }
        }
    }

    document.querySelector('button[data-bs-target="#roomListMenu"]').click();
});

function handlerConnectMessage(data) {
    console.log("client_id: " + data.client_id);
    myClientId = data.client_id;
    stunServerUrl = data.stun_addr;
    turnServerUrl = data.turn_addr;
    requestRoomList();
}

function handlerResultMessage(data) {
    if (data.result_message === "success") {
        showRoomList(data.rooms);
    } else if (data.result_message === "error") {
        alert(data.error_message);
    } else if (data.result_message === "join_success") {
        startStreaming();
    }
}

function handlerOfferMessage(data) {
    console.log("offer received");
    createPeerConnection(data.request_client_id);
    var peerConnection = peerConnectionMap.get(data.request_client_id);

    peerConnection.setRemoteDescription(new RTCSessionDescription(data.offer));
    peerConnection.createAnswer().then(function(answer) {
        peerConnection.setLocalDescription(answer);
        rtcCamSocket.send(JSON.stringify({
            signaling: {
                request_type: 'answer',
                request_client_id: myClientId,
                response_client_id: data.request_client_id,
                answer: answer,
            }
        }));
    });
}

function handlerAnswerMessage(data) {
    console.log("answer received");
    var peerConnection = peerConnectionMap.get(data.request_client_id);
    peerConnection.setRemoteDescription(new RTCSessionDescription(data.answer));
}

function handlerCandidateMessage(data) {
    console.log("candidate received");
    var peerConnection = peerConnectionMap.get(data.request_client_id);
    peerConnection.addIceCandidate(new RTCIceCandidate(data.candidate));
}

function requestRoomList() {
    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'room_list',
        }
    }));
}

function requestRoomJoin(roomId) {
    peerConnectionMap.clear();
    peerVideoStreamMap.clear();
    updateVideoElement();

    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'join_room',
            join_room_id: roomId,
        }
    }));
}

function startStreaming() {
    navigator.mediaDevices.getUserMedia({video: true, audio: true}).then(stream => {
        localMediaStream = stream;

        localVideoElement.srcObject = stream;
        createPeerConnection(0);

        for (let clientId in roomClients) {
            let clientIdInt = parseInt(clientId);
            createPeerConnection(clientIdInt);
            let peerConnection = peerConnectionMap.get(clientIdInt);
            peerConnection.createOffer().then(function(offer) {
                peerConnection.setLocalDescription(offer);
                rtcCamSocket.send(JSON.stringify({
                    signaling: {
                        request_type: 'offer',
                        request_client_id: myClientId,
                        response_client_id: clientIdInt,
                        offer: offer,
                    }
                }));
            });
        }
    }).catch(error => {
        alert("카메라와 오디오를 사용할 수 없습니다. Error: " + error);
    });
}

function createPeerConnection(clientId) {
    var peerConnection = new RTCPeerConnection({
        'iceServers': [
            {'urls': stunServerUrl},
            {'urls': turnServerUrl, 'username': 'test', 'credential': 'test'},
            {'urls': 'turn:kyj9447.iptime.org:50001', 'username': 'test', 'credential': 'test'},
        ]
    });

    localMediaStream.getTracks().forEach(track => peerConnection.addTrack(track, localMediaStream));

    peerConnection.onicecandidate = function(event) {
        if (event.candidate) {
            rtcCamSocket.send(JSON.stringify({
                signaling: {
                    request_type: 'candidate',
                    request_client_id: myClientId,
                    response_client_id: clientId,
                    candidate: event.candidate,
                }
            }));
        }
    }

    peerConnection.ontrack = function(event) {
        peerVideoStreamMap.set(clientId, event.streams[0]);
    }

    peerConnection.oniceconnectionstatechange = function(event) {
        if (['disconnected', 'failed', 'closed'].includes(peerConnection.iceConnectionState)) {
            console.log("peerConnection closed");
            peerConnection.close();
            peerConnectionMap.delete(clientId);
            peerVideoStreamMap.delete(clientId);
            updateVideoElement();
        }
        if (['connected', 'completed'].includes(peerConnection.iceConnectionState)) {
            updateVideoElement();
        }
    }

    if (clientId === 0) {
        localPeerConnection = peerConnection;
    } else {
        peerConnectionMap.set(clientId, peerConnection);
    }


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
            roomClients = room.clients;
            requestRoomJoin(room.id);
            document.querySelector('button[data-bs-target="#roomListMenu"]').click();
        }
        roomElement.style.cursor = "pointer";
        roomElement.innerHTML = htmlString;
        roomListElement.appendChild(roomElement);
    }
}

function updateVideoElement() {
    let peerVideosDiv = document.getElementById('peerVideos');

    while(peerVideosDiv.firstChild) {
        peerVideosDiv.removeChild(peerVideosDiv.firstChild);
    }
    let rowDiv = null;
    let count = 0;

    peerVideoStreamMap.forEach((stream, clientId) => {
        if (count % 2 === 0) {
            rowDiv = document.createElement('div');
            rowDiv.className = "row";
            peerVideosDiv.appendChild(rowDiv);
        }

        let colDiv = document.createElement('div');
        colDiv.className = "col-6";

        let videoElem = document.createElement('video');
        videoElem.srcObject = stream;
        videoElem.autoplay = true;
        videoElem.muted = false;
        videoElem.playsinline = true;
        videoElem.className = "embed-responsive-item";
        videoElem.style.width = `${window.innerWidth / 2}px`;
        videoElem.style.height = `${window.innerWidth / 2}px`;

        colDiv.appendChild(videoElem);
        rowDiv.appendChild(colDiv);

        count++;
    });
}