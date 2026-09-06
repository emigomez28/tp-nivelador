package protocol

type BetBatch struct {
	payload    []byte
	count      int
	maxRecords int
}

func NewBetBatch(maxRecords int) *BetBatch {
	return &BetBatch{payload: []byte{}, count: 0, maxRecords: maxRecords}
}

func (batch *BetBatch) CanAdd(bet Bet) bool {
	if batch.IsEmpty() {
		return true
	}

	maxRecordsNotExceeded := batch.count < batch.maxRecords
	doesNotOverflow := batch.sizeWith(bet) <= _MAX_PAYLOAD_SIZE

	return maxRecordsNotExceeded && doesNotOverflow
}

func (batch *BetBatch) Add(bet Bet) {
	batch.payload = AppendBet(batch.payload, bet)
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

func (batch *BetBatch) Message() Message {
	return NewMessage(MsgBets, batch.payload)
}

func (batch *BetBatch) sizeWith(bet Bet) int {
	return len(batch.payload) + bet.EncodedSize()
}
