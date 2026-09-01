package protocol

const (
	_BET_SEPARATOR      = '\n'
	_BET_SEPARATOR_SIZE = 1
)

type BetBatch struct {
	payload    []byte
	count      int
	maxRecords int
}

func NewBetBatch(maxRecords int) *BetBatch {
	return &BetBatch{payload: []byte{}, count: 0, maxRecords: maxRecords}
}

func (batch *BetBatch) CanAdd(betLine string) bool {
	if batch.IsEmpty() {
		return true
	}

	maxRecordsNotExceeded := batch.count < batch.maxRecords
	doesNotOverflow := batch.sizeWith(betLine) <= _MAX_PAYLOAD_SIZE

	return maxRecordsNotExceeded && doesNotOverflow
}

func (batch *BetBatch) Add(betLine string) {
	if batch.count > 0 {
		batch.payload = append(batch.payload, _BET_SEPARATOR)
	}

	for _, betLineByte := range []byte(betLine) {
		batch.payload = append(batch.payload, betLineByte)
	}

	batch.count++
}

func (batch *BetBatch) Count() int {
	return batch.count
}

func (batch *BetBatch) IsEmpty() bool {
	return batch.count == 0
}

func (batch *BetBatch) Reset() {
	batch.payload = batch.payload[:0]
	batch.count = 0
}

func (batch *BetBatch) Message() *Message {
	return NewMessage(MsgBets, batch.payload)
}

func (batch *BetBatch) sizeWith(betLine string) int {
	size := len(batch.payload) + len(betLine) + _BET_SEPARATOR_SIZE
	return size
}
