Nanomsg Queue
=============

A playground repo, where I wanted to play around with nanomsg in golang.

## How to use
You have producers and consumers.

### Producers
Adds an item to be worked on. The command will block until a consumer takes and
acks the work.
```
## work value is the number, in seconds, you want the task to take
go run main.go todo '{"work":1}'
```

### Consumers
Consumers need a url to get the work to do from. I've harded code three
options, `tcp://localhost:40897`, `tcp://localhost:40898`,
`tcp://localhost:40899`.

```
go run main.go pop "tcp://localhost:40899"
```
