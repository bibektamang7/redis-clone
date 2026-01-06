<img width="1593" height="767" alt="Image" src="https://github.com/user-attachments/assets/7b2fc2f9-e020-4dc6-bb49-52a2205fd7a3" />

### Request Structure

```go
// the structure of redis request is similar to HTTP request format
// *3\r\n$3\r\nset\r\n$5\r\nadmin\r\n$5bibek
*3
$3
set
$5
admin
$5
bibek

// * -> represents array
// $ -> represents string
// NOTE: MAKE SURE ITS CORRECT
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

$5\r\nbibek\r\n

$5
bibek

```

RESP 3

```json
> HELLO 3
1# "server" => "redis"
2# "version" => "6.0.0"
3# "proto" => (integer) 3
4# "id" => (integer) 10
5# "mode" => "standalone"
6# "role" => "master"
7# "modules" => (empty array)
```
