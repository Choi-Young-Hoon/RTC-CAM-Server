// RTCCAM 서버랑 통신 Websocket
var rtcCamSocket = null;
var rtcCamWSServerUrl = "{{.WebSocketURL}}";

// WebRTC
var currentClientId = 0;
var peerConnectionMap = new Map();
var peerVideoStreamMap = new Map();
var iceServers = [];

// Local 영상
var localMediaStream = null;

window.onload = function() {
    initRTCCamSocket();
};

// Android 에서 video 전체화면 재생 활성화
function localVideoFullScreen() {
    localFullScreen(true);
}

// Android 에서 video 전체화면 재생 해제
function localVideoExitFullScreen() {
    localFullScreen(false);
}

function localFullScreen(isFullScreen) {
    var videoContainer = document.getElementById('videoContainer');
    if (!isFullScreen) {
        videoContainer.style.display = "none";
    } else {
        videoContainer.style.display = "flex";
    }
}

function initRTCCamSocket() {
    rtcCamSocket = new WebSocket(rtcCamWSServerUrl);
    rtcCamSocket.onopen = function () {
        console.log("WebSocket opened");
        requestRoomList();

        var currentUrl = window.location.href;
        if (currentUrl.includes("/room?join_room=")) {
            let params = new URLSearchParams(window.location.search);
            let joinRoomId = params.get('join_room');
            let authToken = params.get('auth_token');
            console.log("join_room: " + joinRoomId + ", auth_token: " + authToken);
            requestRoomJoin(joinRoomId, authToken);
        }
    }

    rtcCamSocket.onerror = function (event) {
        alert("rtccam 서버와 통신할 수 없습니다.");
        rtcCamSocket = new WebSocket(rtcCamWSServerUrl);
    }

    rtcCamSocket.onclose = function (event) {
        console.log("WebSocket closed");
    }

    rtcCamSocket.onmessage = function (event) {
        const data = JSON.parse(event.data);
        console.log(data);

        if (data.room_id !== undefined) {
            handlerCreateRoom(data);
        } else if (data.auth_token !== undefined) { // auth token
            handlerAuthToken(data);
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

function handlerAuthToken(data) {
    moveRoom(data.room_info.id, data.auth_token);

}

function handlerCreateRoom(data) {
    moveRoom(data.room_id, data.auth_token);
}

function handlerResultMessage(data) {
    if (data.result_message === "success") {
        showRoomList(data.rooms);
    } else if (data.result_message === "error") {
        alert(data.error_message);
    } else if (data.result_message === "join_success") {
        iceServers = data.ice_servers;
        currentClientId = data.client_id;
        startStreaming(data.room_info);
    } else if (data.result_message === "leave_client") {
        peerClose(data.client_id);
    }
}

function handlerOfferMessage(data) {
    console.log("offer received");
    createPeerConnection(data.request_client_id);
    var peerConnection = peerConnectionMap.get(data.request_client_id);

    peerConnection.setRemoteDescription(new RTCSessionDescription(data.offer));
    peerConnection.createAnswer().then(function(answer) {
        peerConnection.setLocalDescription(answer);
        requestAnswer(data.request_client_id, answer);
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
    console.log("candidate data peerConnection: " + peerConnection);
    console.log("candidate data candidate: " + data.candidate);
    peerConnection.addIceCandidate(new RTCIceCandidate(data.candidate));
}

function clearPeerMap() {
    peerConnectionMap.clear();
    peerVideoStreamMap.clear();
    updateVideoElement();
}

function requestCreateRoom(roomTitle, isPassword, roomPassword) {
    rtcCamSocket.send(JSON.stringify({
        create_room: {
            title: roomTitle,
            password: isPassword ? roomPassword : '',
        },
    }));
}

function requestAuthToken(roomId, password) {
    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'auth_token',
            password: password,
            room_id: roomId,
        },
    }));
}

function requestRoomList() {
    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'room_list',
        }
    }));
}

function requestRoomJoin(roomId, auth_token) {
    //clearPeerMap();
    console.log("requestRoomJoin: " + roomId + ", " + auth_token);

    rtcCamSocket.send(JSON.stringify({
        room: {
            request_type: 'join_room',
            auth_token: auth_token,
            room_id: Number(roomId),
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

function requestOffer(clientId, offer) {
    rtcCamSocket.send(JSON.stringify({
        signaling: {
            request_type: 'offer',
            request_client_id: currentClientId,
            response_client_id: clientId,
            offer: offer,
        }
    }));
}

function requestAnswer(clientId, answer) {
    rtcCamSocket.send(JSON.stringify({
        signaling: {
            request_type: 'answer',
            request_client_id: currentClientId,
            response_client_id: clientId,
            answer: answer,
        }
    }));
}

function requestCandidate(clientId, candidate) {
    rtcCamSocket.send(JSON.stringify({
        signaling: {
            request_type: 'candidate',
            request_client_id: currentClientId,
            response_client_id: clientId,
            candidate: candidate,
        }
    }));
}

function startStreaming(roomInfo) {
    navigator.mediaDevices.getUserMedia({video: true, audio: true}).then(stream => {
        localMediaStream = stream;
        updateVideoElement();

        for (let clientId in roomInfo.clients) {
            if (parseInt(clientId) === currentClientId) {
                continue;
            }

            let clientIdInt = parseInt(clientId);
            createPeerConnection(clientIdInt);
            let peerConnection = peerConnectionMap.get(clientIdInt);
            peerConnection.createOffer().then(function(offer) {
                peerConnection.setLocalDescription(offer);
                requestOffer(clientIdInt, offer);
            });
        }

    }).catch(error => {
        alert("카메라와 오디오를 사용할 수 없습니다. Error: " + error);
    })
}

function peerClose(clientId) {
    console.log("[peerClose] peerConnection closed");
    peerConnectionMap.delete(clientId);
    peerVideoStreamMap.delete(clientId);
    updateVideoElement();
}

function createPeerConnection(clientId) {
    var peerConnection = new RTCPeerConnection({
        "iceServers": iceServers
    });

    localMediaStream.getTracks().forEach(track => peerConnection.addTrack(track, localMediaStream));

    peerConnection.onicecandidate = function(event) {
        if (event.candidate) {
            requestCandidate(clientId, event.candidate);
        }
    }

    peerConnection.onnegotiationneeded = function() {
        peerConnection.createOffer().then(function(offer) {
            return peerConnection.setLocalDescription(offer);
        })
            .then(function() {
                requestOffer(clientId, peerConnection.localDescription)
            })
            .catch(function(error) {
                console.error("Error creating offer: ", error);
            });
    };
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

    peerConnectionMap.set(clientId, peerConnection);
}

function createVideoElement(stream, className, videoName) {
    var div = document.createElement('div');
    div.className = className;

    // 새로운 video 요소를 생성합니다.
    var video = document.createElement('video');
    video.srcObject = stream;
    video.autoplay = true;
    video.controls = false;
    video.width = "auto";
    video.height = "auto";

    // 새로운 div 요소를 생성합니다.
    var divName = document.createElement('div');
    divName.className = "participant-name";
    divName.textContent = videoName;

    // div 요소에 video와 divName 요소를 추가합니다.
    div.appendChild(video);
    div.appendChild(divName);

    return div;
}

function updateVideoElement() {
    var videoSection = document.getElementById('videoSection');

    while (videoSection.firstChild) {
        videoSection.removeChild(videoSection.firstChild);
    }

    // 로컬 비디오 표시
    var localVideoDiv = createVideoElement(localMediaStream, "my-video-container main-participant", "나");
    videoSection.appendChild(localVideoDiv);

    // Peer 들의 비디오 표시
    var loopCount = 0;
    peerConnectionMap.forEach(function(peerConnection, clientId) {
        var video = createVideoElement(peerVideoStreamMap.get(clientId), "video-container peer-participant", "참여자 " + ++loopCount);
        videoSection.appendChild(video);
    });
}

function showRoomList(roomList) {
    var roomTable = document.getElementById('roomTable');

    while (roomTable.rows.length > 1) {
        roomTable.deleteRow(1);
    }

    for (let roomId in roomList.rooms) {
        let room = roomList.rooms[roomId];
        let row = roomTable.insertRow(1);
        let cell1 = row.insertCell(0);
        let cell2 = row.insertCell(1);
        let cell3 = row.insertCell(2);

        cell1.innerHTML = room.id;
        cell2.innerHTML = room.title;

        let clientCount = Object.keys(room.clients).length === 0 ? 0 : Object.keys(room.clients).length;
        cell3.innerHTML = clientCount;

        row.onclick = function() {
            var password = "";
            if (room.is_password) {
                password = prompt("비밀번호를 입력해주세요.");
            }
            requestAuthToken(room.id, password);
        }

    }

    initRoom = true;
}






//////////////////////////////////////// 왼쪽 메뉴 ////////////////////////////////////////
function moveRoom(roomId, authToken) {
    var isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
    if (isMobile) {
        //window.location.href = "/room_mobile?join_room=" + roomId.toString() + "&auth_token=" + authToken;
        window.location.href = "/room?join_room=" + roomId.toString() + "&auth_token=" + authToken;
    } else {
        window.location.href = "/room?join_room=" + roomId.toString() + "&auth_token=" + authToken;
    }
}

function openMenu() {
    if (document.getElementById('mySidenav').style.width !== "0px") {
        closeMenu();
    } else {
        document.getElementById("mySidenav").style.width = "300px";
        document.getElementById("rtccamOverlay").style.width = "100%";
    }
}

function closeMenu() {
    document.getElementById("mySidenav").style.width = "0";
    document.getElementById("rtccamOverlay").style.width = "0";
}


//////////////////////////////////////// 방 생성 모달 ////////////////////////////////////////
var createRoomModal = document.getElementById('createRoomModal');

function showCreateRoomModal() {
    createRoomModal.style.display = "block";
}

function hideCreateRoomModal() {
    createRoomModal.style.display = "none";
}

document.getElementById('usePassword').addEventListener('change', function() {
    var passwordField = document.getElementById('passwordField');
    var roomPassword = document.getElementById('roomPassword');
    if (this.checked) {
        passwordField.classList.remove('hidden');
        roomPassword.disabled = false;
    } else {
        passwordField.classList.add('hidden');
        roomPassword.disabled = true;
    }
});

function createRoom() {
    var roomTitle = document.getElementById('roomTitle').value;
    var isPassword = document.getElementById('usePassword').checked;
    var roomPassword = document.getElementById('roomPassword').value;

    if (roomTitle === "") {
        alert("방 제목을 입력해주세요.");
        return;
    }

    requestCreateRoom(roomTitle, isPassword, roomPassword);
    hideCreateRoomModal();
}
