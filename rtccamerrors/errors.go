package rtccamerrors

import "errors"

var ErrorClientNotFound = errors.New("클라이언트를 찾을 수 없습니다.")

var ErrorRoomNotFound = errors.New("방을 찾을 수 없습니다.")
var ErrorRoomIsFull = errors.New("방이 꽉 찼습니다.")

var ErrorRequestTypeError = errors.New("요청 타입이 올바르지 않습니다.")
var ErrorTitleIsEmpty = errors.New("방 제목이 비어있습니다.")
var ErrorInvalidMaxClientCount = errors.New("최대 인원 수가 올바르지 않습니다.")
var ErrorInvalidPassword = errors.New("비밀번호가 올바르지 않습니다.")
var ErrorInvalidAuthToken = errors.New("잘못된 접근입니다.")
