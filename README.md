# AdHocTX-Go
Transactions on multiple data types over ad-hoc infrastructure

# Theory

The purpose of AdHocTX is to provide a portable ID that represents some version of data, and 
build upon that primitive to create performant, decentralized, geo-scalable, and ACID-compliant databases
through different mediums from in-memory, to disk, to cloud storage.

The only underlying assumption is that there is no need for a single authorative version, and no version
is considered the true version. The only guarantee that people need is that no one is going to be trusting
a version that does not contain everything that yours contains. The owner of the data is able to
produce new versions freely, and must only ensure that any changes do that version do not effect any
other versions.

This means when you want to read data, you can use whatever version you have lying around. You can also request
a new version from the data owner, and make whatever changes you like to that data. However, before considering 
your new version as authorative, you just have to make sure that it contains all the data in everyone elses
authorative version. You do this by just asking them to consider yours authorative.

When you ask someone to consider yours authorative, it doesn't mean they are going to publish it. It just means
that, if they want to start publishing a new version, they will ask you first. If you end up not publishing your
version, you can just go and tell everyone to forget about it, just to save them the hassle later.

The speed comes from being able to send your peers your version before you even start changing it. They will
make sure that they don't publish anything incompatible with it, so you can publish as soon as you've received 
confirmations from your peers. Then, in a hurry, tell the other peers you've actually published it, in case 
they were waiting on it.

In simpler terms, publishing a new version in CA guarantees that someone in the EU isn't going to publish a conflicting
version. They will, however, have to wait for you to inform them, before they can start publishing.

If you forget to inform them (IE, if the node dies), they will continue to be able to prepare new conflicting version, 
but none of them will be able to publish it. You may have already published a version to someone before you died, and
that's fine. However, when you finally come back up, 

This is eventual consistency. Strong consistency is a single flip of the switch. Whereas before, no one could publish a
**new** version before you finished with your version, strong consistency means that no one can publish **any** version
before you're finished - though they can still prepare the response.

# Flexibility

The use of such a simple primitive allows for high flexibility of CAP theorem.

If you have a high-read KV store acting as a cache, writing to it would only have the delay of your furthest node. 
The latency peaks would only occur when 2 nodes are both wanting to publish a version. Suppose 2 nodes concurrently 
send a reserve request. 20ms later, both will respond with an ACK, and will include a stream of their version and a 
FIN because they know there is a version conflict. Upon receiving, usually they would apply the changes and publish
different versions - their own versions, and the diverged versions would be cleaned up later. However, because they
edited the same key, they reference each others "listeners" file in the data store, and use the child version within to
determine which commit fails and which commit succeeds. If one node received the reserve request before it started a
new version, it would send it's usual ACK, but now know to send that node a stream with it's own reserve request. In 
summary, a **merge constraint violations have 1 roundtrip to resolve**, identical to a no-violation situation. The latency
is still high, but it is required: the values written may have been derived from the values read.

# Isolation / Consistency

There are 3 ways to read a version, from weakest to strongest:
- READ_LATEST, which reads the most-recently received reserved version that may be thrown away
- READ_PUBLISHED, which reads the published version, **default**
- READ_STRONG, which returns the published version as long as there are no reserved versions

There are 2 modes for publishing a version:
- PUB_RESERVE, used when the changes may violate constraints, will reserve if needed **default**
- PUB_FORCE, publishes a version that may violate a constraint when merged with peers and be lost

There are 3 generally-obeyed modes of consistency, a mode or higher will at least be obeyed
- SAVE_MEM, to ensure the data is saved in memory
- SAVE_DISK, to ensure the data has been saved and synced on disk **default**
- SAVE_COPIES, to ensure the data has been replicated by at least 4 nodes, in 2 geo's if available

Finally, there are 3 levels of urgency, which is used to determine allocated workers for larger queries
- SPEED_PAIR, which will not engage more than 2 nodes
- SPEED_CLUSTER, which will utilize the free-time of the cluster **default**
- SPEED_SCALE, which will allocate more resources in whatever ways are configured
