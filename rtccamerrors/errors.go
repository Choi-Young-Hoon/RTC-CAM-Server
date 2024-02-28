package rtccamerrors

import "errors"

var ErrorClientNotFound = errors.New("클라이언트를 찾을 수 없습니다.")

var ErrorRoomNotFound = errors.New("방을 찾을 수 없습니다.")

var ErrorRequestTypeError = errors.New("요청 타입이 올바르지 않습니다.")
var ErrorTitleIsEmpty = errors.New("방 제목이 비어있습니다.")

var ErrorNotFoundCreateRoomWaitItem = errors.New("방 생성 대기열에서 아이템을 찾을 수 없습니다.")