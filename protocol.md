1. **Event Type**
   - Offset: 0 bytes
   - Length: 1 byte
   - Byte 0:
     - `0x00` - send data
     - `0x01` - subscribe
     - `0x02` - unsubscribe

2. **Producer & Stream identifier**
   - Offset: 1 byte
   - Length: 4 bytes
   - Byte 1-3: Producer ID
   - Byte 4: Stream ID
