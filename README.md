# 🚀 Mini Kafka (Go)

A Kafka-inspired distributed messaging system built from scratch in Go, implementing core concepts like log-based storage, partitioning, consumer groups, offset management, and TCP-based communication.

---

## 🧠 Overview

This project is a simplified implementation of a distributed log-based messaging system inspired by Apache Kafka.

It focuses on understanding how real-world streaming systems work internally, including:

* Append-only log storage
* Partition-based scalability
* Consumer group coordination
* Offset tracking & recovery
* TCP-based custom protocol
* Basic replication concepts

---

## 🏗️ Architecture

```
Producer → Broker → Storage (Log Files)
                  → Consumer Groups
                  → Offset Manager
                  → TCP Server
```

---

## ⚙️ Features

### ✅ Core Features

* 📦 **Persistent Log Storage**

    * Append-only file-based storage
    * Offset-based message retrieval

* 🧩 **Topics & Partitions**

    * Multi-partition support
    * Hash-based partition routing

* 👥 **Consumer Groups**

    * Partition assignment
    * Load balancing
    * Per-group offset tracking

* 💾 **Offset Management**

    * Persistent offsets (disk-based)
    * Crash recovery support

* 🌐 **Custom TCP Protocol**

    * Length-prefixed messaging
    * Persistent connections
    * Multiple requests per connection

---

### 🚀 Advanced Features

* 🔁 **Leader-Follower Replication (Basic)**

    * Asynchronous replication
    * Multi-broker simulation

* 🔄 **Rebalancing**

    * Partition reassignment on consumer join

---

## 🧠 Key Concepts Implemented

| Concept         | Description                             |
| --------------- | --------------------------------------- |
| Log Storage     | Messages stored sequentially in files   |
| Offset          | Byte position used for reading messages |
| Partitioning    | Horizontal scaling mechanism            |
| Consumer Groups | Parallel consumption model              |
| Metadata        | Topic & partition information           |
| TCP Protocol    | Efficient network communication         |

---

## 📁 Project Structure

```
mini-kafka/
├── broker/
│   ├── broker.go
│   ├── topic.go
│   ├── partition.go
│   ├── consumer_group.go
├── storage/
│   ├── log.go
│   ├── offset.go
├── protocol/
├── main.go
```

---

## 🚀 Getting Started

### 1️⃣ Clone Repository

```bash
git clone https://github.com/your-username/mini-kafka.git
cd mini-kafka
```

---

### 2️⃣ Run Broker

```bash
go run main.go
```

Server starts on:

```
localhost:9092
```

---

## 🧪 Example Requests (TCP)

### 📤 Produce

```json
{
  "type": "produce",
  "topic": "orders",
  "key": "user1",
  "message": "order1"
}
```

---

### 📥 Consume

```json
{
  "type": "consume",
  "topic": "orders",
  "group": "g1",
  "consumer_id": "c1"
}
```

---

## 🔌 TCP Protocol

Custom protocol:

```
[4 bytes length][JSON payload]
```

Example:

```
\x00\x00\x00\x2A {"type":"produce",...}
```

---

## 🧠 How It Works

1. Producer sends message via TCP
2. Broker routes message to partition
3. Message is appended to log file
4. Consumer group reads using offsets
5. Offsets are persisted for recovery

---

## ⚠️ Limitations

This is a learning-focused implementation and does not include:

* ❌ Distributed consensus (Raft / KRaft)
* ❌ Strong consistency guarantees
* ❌ Exactly-once semantics
* ❌ High-performance optimizations (batching, zero-copy)

---

## 📈 Future Improvements

* Leader election using Raft
* Heartbeat-based failure detection
* Metadata service
* Batch processing & compression
* CLI tools (producer/consumer)
* Monitoring & metrics

---

## 🧠 Learnings

Through this project, I explored:

* Designing append-only storage systems
* Handling concurrency and synchronization
* Building TCP-based protocols
* Understanding distributed system fundamentals
* Implementing fault tolerance concepts

---

## 💡 Inspiration

Inspired by Apache Kafka internals and distributed event streaming systems.

---

## ⭐ If you found this useful

Give it a ⭐ on GitHub!
