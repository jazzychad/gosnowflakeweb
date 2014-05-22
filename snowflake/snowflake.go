package snowflake

import (
	"time"
)

var workerIdBits uint64 = 5
var datacenterIdBits uint64 = 5
var sequenceBits uint64 = 12

var workerIdShift uint64 = sequenceBits
var datacenterIdShift uint64 = sequenceBits + workerIdBits
var timestampLeftShift uint64 = sequenceBits + workerIdBits + datacenterIdBits
var sequenceMask int64 = int64(1 << sequenceBits) - 1

var twepoch int64 = 1288834974657

type snowflake struct {
	workerId uint64
	datacenterId uint64
	sequenceNumber uint64
	lastTimestamp int64
}

func NewSnowflake(datacenterId uint64, workerId uint64) *snowflake {
	s := new(snowflake)
	s.datacenterId = datacenterId
	s.workerId = workerId
	s.sequenceNumber = 0
	s.lastTimestamp = 0
	return s
}

func (this *snowflake) NextId() uint64 {

	timestamp := timeGen()

	if (timestamp < this.lastTimestamp) {
		panic("clock is going backward!")
	}

	if (this.lastTimestamp == timestamp) {
		this.sequenceNumber = uint64((this.sequenceNumber + 1) & uint64(sequenceMask))
		if (this.sequenceNumber == 0) {
			timestamp = tilNextMillis(this.lastTimestamp)
		}
	} else {
		this.sequenceNumber = 0
	}

	this.lastTimestamp = timestamp

	return uint64(((timestamp - twepoch) << timestampLeftShift)) | 
		(this.datacenterId << datacenterIdShift) |
		(this.workerId << workerIdShift) |
		this.sequenceNumber
}

func (this *snowflake) WorkerId() uint64 {
	return this.workerId
}


//// package private

func timeGen() int64 {
	nanos := time.Now().UnixNano()
	millis := nanos / 1000000
	return millis
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}
