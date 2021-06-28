# Telemetry challenge project

The challenge project consists of building a telemetry visualisation solution for two different satellites. Each satellite has its own encoding described below. The solution must allow users to explore telemetry data from different time ranges.

The ground station provider offers a separate TCP server for each satellite. Packets are sent over this connection without additional framing. The TCP server will accept any number of connections but only the last connected client for each satellite will receive telemetry. The servers will ignore any bytes sent to them.

Candidates are free to use any available open-source software as long as its license allows free commercial use. There is no restriction about the environment the binary must run on.

## Objective

Create a solution that you consider as close to production-ready as possible. Focus on quality rather than quantity.

## Expected output

Candidates are expected to: 
- Provide the source code (link to a public repository or in a compressed folder)
- Include clear instructions on how to execute this code
- Prepare a short presentation (10min approx) about the solution used, the approach taken, the main challenges and next-steps

## Time-frame

To be agreed with each candidate depending on availability and experience. Acceptable values range from one working day to five.

## Technical details

### Telemetry points

Each point contains:
1. Unix timestamp - int64 - from the moment the point was generated
2. Telemetry ID - uint16 - an identifier for the metric (example: a voltage from a specific sensor might be id 1, a reading from a different sensor might be id 2, etc)
3. Value - float32 - value the sensor read

- **String**
  Packets encoded in string format are encoded using the UTF-8 standard. Messages are wrapped in `[]` and the different values are separated by `:`. Values appear in the order described above.
  Example: `[1604614491:1:2.000000]`
- **Binary**
  Packets encoded in binary format use a little-endian encoding. They have a header of 4 bytes with always the same value: `00 01 02 03`. The following 8 bytes represent the timestamp, the next 2 bytes represent the telemetry ID, and the final 4 bytes correspond to the telemetry value.
  Example:
  ```
  00000000  00 01 02 03 ff 79 a4 5f  00 00 00 00 01 00 00 00  |.....y._........|
  00000010  00 40                                             |.@|
  ```

### Binary usage

The provided binary will emulate a ground station provider. When running the command `generate`, it will accept connections on the ports indicated with the arguments. By default port 8000 and 8001 will be used to accept connection for the string-encoding and binary-encoding respectively. Use the optional flags `--portString` and `--portBinary` to change the default values.

The following examples are equivalent:
```
./telemetry generate
./telemetry generate --portString 8000 --portBinary 8001
```
