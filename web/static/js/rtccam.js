var rtcCamSocket = null;
var rtcCamWSServerUrl = 'ws://' + window.location.hostname + ':40001/rtccam';

var myClientId = 0;
var stunServerUrl = null;
var turnServerUrl = null;

var roomClients = [];
var peerConnectionMap = new Map();
var peerVideoStreamMap = new Map();

var localMediaStream = null;
var localVideoElement = document.getElementById('localVideo');

window.onload = function() {
    var loadingDiv = document.getElementById('loadingDiv');
    var overlay = document.getElementById('overlay');
    loadingDiv.style.display = 'block';
    overlay.style.display = 'block';

    setTimeout(function() {
        loadingDiv.style.display = 'none';
        overlay.style.display = 'none';
    }, 1000);
};

document.addEventListener('DOMContentLoaded', function () {


    document.querySelector('button[data-bs-target="#roomListMenu"]').click();

    initRTCCamSocket();

    navigator.mediaDevices.getUserMedia({video: true, audio: true}).then(stream => {
        localMediaStream = stream;

        localVideoElement.srcObject = stream;
        localVideoElement.onloadedmetadata = function(e) {
            localVideoElement.play();
        }
    }).catch(error => {
        alert("카메라와 오디오를 사용할 수 없습니다. Error: " + error);
    });
});

localVideoElement.addEventListener('click', function() {
    localVideoElement.requestPictureInPicture().catch(error => {
        console.error(error);
    });
});

document.addEventListener('visibilitychange', function() {
    try {
        if (document.hidden && 'pictureInPictureEnabled' in document && isLocalVideoPIP) {
            localVideoElement.requestPictureInPicture().catch(error => {
                console.error(error);
            });
        }
    } catch (error) {
        console.error(error);
    }
});

function initRTCCamSocket() {
    rtcCamSocket = new WebSocket(rtcCamWSServerUrl);
    rtcCamSocket.onopen = function () {
        console.log("WebSocket opened");
    }

    rtcCamSocket.onerror = function (event) {
        alert("rtccam 서버와 통신할 수 없습니다.");
        rtcCamSocket = new WebSocket(rtcCamWSServerUrl);
    }

    rtcCamSocket.onmessage = function (event) {
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
}

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
    } else if (data.result_message === "leave_client") {
        peerClose(data.leave_client_id);
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

function clearPeerMap() {
    peerConnectionMap.clear();
    peerVideoStreamMap.clear();
    updateVideoElement();
}

async function requestCreateRoom(roomTitle, isPassword, roomPassword) {
    roomClients = [];
    clearPeerMap();

    await requestRoomLeave(myClientId);
    await rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'create_room',
            title: roomTitle,
            password: isPassword ? roomPassword : '',
        },
    }));
}

async function requestRoomList() {
    await rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'room_list',
        }
    }));
}

function requestRoomJoin(roomId, password) {
    clearPeerMap();

    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'join_room',
            password: password,
            join_room_id: roomId,
        }
    }));
}

async function requestRoomLeave(clientId) {
    await rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'leave_room',
        }
    }));
}

function startStreaming() {
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

    localVideoElement.requestPictureInPicture();
    document.querySelector('button[data-bs-target="#roomListMenu"]').click();
}

function peerClose(clientId) {
    console.log("[peerClose] peerConnection closed");
    peerConnectionMap.delete(clientId);
    peerVideoStreamMap.delete(clientId);
    updateVideoElement();
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
            peerClose(clientId);
        }
        if (['connected', 'completed'].includes(peerConnection.iceConnectionState)) {
            updateVideoElement();
        }
    }

    if (clientId === 0) {

    } else {
        peerConnectionMap.set(clientId, peerConnection);
    }


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
            if (room.is_password === true) {
                requestRoomJoin(room.id, password=prompt("비밀번호를 입력하세요"));
            } else {
                requestRoomJoin(room.id, '');
            }
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

    if (count % 2 === 0) {
        rowDiv = document.createElement('div');
        rowDiv.className = "row";
        peerVideosDiv.appendChild(rowDiv);
    }

    peerVideoStreamMap.forEach((stream, clientId) => {
      /*  if (count % 2 === 0) {
            rowDiv = document.createElement('div');
            rowDiv.className = "row";
            peerVideosDiv.appendChild(rowDiv);
        }
*/
        let colDiv = document.createElement('div');
        colDiv.className = "col-sm-6 col-md-4 col-lg-3 col-xl-2";

        let videoElem = document.createElement('video');
        videoElem.srcObject = stream;
        videoElem.autoplay = true;
        videoElem.muted = false;
        videoElem.playsinline = true;
        videoElem.style.width = "100%";
        videoElem.style.height = "100%";
        videoElem.addEventListener('click', function() {
            if (videoElem.requestFullscreen) {
                videoElem.requestFullscreen();
            } else if (videoElem.mozRequestFullScreen) { // Firefox
                videoElem.mozRequestFullScreen();
            } else if (videoElem.webkitRequestFullscreen) { // Chrome, Safari and Opera
                videoElem.webkitRequestFullscreen();
            } else if (videoElem.msRequestFullscreen) { // IE/Edge
                videoElem.msRequestFullscreen();
            }
        });

        colDiv.appendChild(videoElem);
        rowDiv.appendChild(colDiv);
        peerVideosDiv.appendChild(rowDiv);
        count++;
    });
}


let createRoomModal = new bootstrap.Modal(document.getElementById('createRoomModal'));
////////////////////////////////////////////////////////////////////////////////////////////////////// Modal
document.getElementById('createRoomButton').addEventListener('click', showCreateRoomDialog);
function showCreateRoomDialog() {
    createRoomModal.show();
}

document.getElementById('createRoomButton').addEventListener('click', function() {
    var roomTitle = document.getElementById('roomTitle').value;
    var isPassword = document.getElementById('isPassword').checked;
    var roomPassword = document.getElementById('roomPassword').value;

    console.log('Room Title:', roomTitle);
    console.log('isPassword:', isPassword);
    console.log('Room Password:', roomPassword);

    requestCreateRoom(roomTitle, isPassword, roomPassword);

    createRoomModal.hide();
});

document.getElementById('isPassword').addEventListener('change', function() {
    document.getElementById('roomPassword').disabled = !this.checked;
});
