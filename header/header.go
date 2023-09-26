package headers

func Encode(event uint8, producerID int32) []byte {
	buffer := make([]byte, 5)

	buffer[0] = event
	buffer[1] = byte(producerID >> 24)
	buffer[2] = byte(producerID >> 16)
	buffer[3] = byte(producerID >> 8)
	buffer[4] = byte(producerID)

	return buffer
}

func Decode(buffer []byte) (uint8, int32) {
	event := buffer[0]

	producerID := int32(buffer[1])<<24 |
		int32(buffer[2])<<16 |
		int32(buffer[3])<<8 |
		int32(buffer[4])

	return event, producerID
}
