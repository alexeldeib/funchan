# funchan

> fun (with) chan(nels)

A repository to play with data structures in Go, particularly those based on channels.

- broker: a simple pub/sub broker allowing delivery of a single message to many subscribers. 
- queue: an extremely simple push/pop/peek queue based on slice
- set: a map-based set with has/insert/delete
- priorityqueue: defines a Heapable interface and uses it to implement a generic priority queue.
- workqueue: a priority queue based on time which pops the root item into an out channel when current time is greater than the root time.
- broadcaster: similar to broker, allows broadcasting messages to many consumers. No unsubscription, no blocking writes.
- notifier: similar to broadcaster, allows many consumers to block on a single notification. Broadcaster allows for sending additional data with the signal, notifier does not.