# Anchor

A high-performance, distributed key-value database built from the ground up with a focus on reliability, consistency, and scalability.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue.svg)]()
![Status](https://img.shields.io/badge/status-under%20development-orange)

> ⚠️ **Note**: Anchor is currently under active development. The API and features are subject to change. Not recommended for production use at this time.

## Overview

Anchor is a distributed database system that provides strong consistency guarantees, automatic failover, and horizontal scalability. It combines battle-tested data structures with modern distributed systems algorithms to deliver a robust storage solution for mission-critical applications.

### Key Highlights

- 🚀 **High Performance**: Optimized storage engines with minimal write amplification
- 🔒 **ACID Transactions**: Full support for distributed transactions with strong consistency
- 📊 **Multiple Storage Engines**: Choose between B-Tree and LSM-Tree based on your workload
- 🌐 **Horizontal Scalability**: Add nodes dynamically to scale read and write throughput
- 💪 **Fault Tolerant**: Automatic failure detection and recovery with zero data loss
- ⚡ **Low Latency**: Optimized data structures and caching for sub-millisecond reads
- 🔄 **Flexible Consistency**: Support for both strong and eventual consistency models

---

## Features

### Storage Engine

#### B-Tree Storage Engine
- **Disk-optimized B-Tree implementation** with configurable order
- **Page-based storage** with efficient space utilization
- **Slotted page format** supporting variable-length records
- **Sibling pointers** for efficient range scans
- **Prefix compression** to reduce storage overhead
- **Copy-on-Write (CoW) variant** for MVCC and snapshot isolation
- **Bulk loading optimization** for faster initial data loading
- **Automatic page defragmentation** to maintain performance

#### LSM-Tree Storage Engine
- **Write-optimized Log-Structured Merge-Tree** architecture
- **In-memory MemTable** with fast skiplist implementation
- **Immutable SSTables** with block-based organization
- **Bloom filters** for efficient negative lookups
- **Leveled and Tiered compaction** strategies
- **Block compression** (Snappy, LZ4, Zstd)
- **Tombstone-based deletion** with efficient garbage collection
- **Configurable write and read amplification** trade-offs

#### Storage Features
- **Custom binary page format** with checksums for data integrity
- **Magic numbers and version headers** for format validation
- **Zero-copy serialization** for performance
- **Page-level encryption** support (AES-256)
- **Configurable page sizes** (4KB to 64KB)
- **Smart buffer pool** with LRU eviction policy
- **Direct I/O support** to bypass OS cache when needed

---

### Transaction Processing

#### ACID Guarantees
- **Atomicity**: All-or-nothing transaction execution
- **Consistency**: Maintain database invariants across transactions
- **Isolation**: Serializable snapshot isolation (SSI)
- **Durability**: Write-ahead logging with configurable sync modes

#### Transaction Features
- **Write-Ahead Log (WAL)** with LSN-based ordering
- **Group commit** optimization for higher throughput
- **ARIES-style recovery** (Analysis, Redo, Undo phases)
- **Automatic checkpointing** to limit recovery time
- **Savepoints** for partial rollback within transactions
- **Distributed transactions** with two-phase commit (2PC)
- **Optimistic concurrency control** to minimize lock contention
- **Deadlock detection and prevention** with timeout-based abortion
- **Multi-version concurrency control (MVCC)** for read isolation
- **Snapshot isolation** for consistent point-in-time reads

---

### Distributed System Architecture

#### Cluster Management
- **Dynamic cluster membership** with auto-discovery
- **Gossip-based membership protocol** (SWIM)
- **Phi-accrual failure detection** with adaptive thresholds
- **Automatic node replacement** on failures
- **Rolling upgrades** with zero downtime
- **Cluster metadata service** for topology management
- **Health monitoring** with configurable health checks

#### Consensus and Replication
- **Raft consensus algorithm** for leader election
- **Strong consistency** with linearizable reads and writes
- **Configurable replication factor** (default: 3)
- **Quorum-based reads and writes** for tunable consistency
- **Log replication** with pipelining for throughput
- **Log compaction** via snapshotting
- **Single-server membership changes** for safe reconfiguration
- **Leader lease** for optimized read performance
- **Pre-vote optimization** to prevent election disruption

#### Data Distribution
- **Consistent hashing** for balanced data distribution
- **Virtual nodes** for fine-grained load balancing
- **Range-based partitioning** option for locality
- **Automatic rebalancing** when nodes join or leave
- **Partition placement optimization** based on load
- **Cross-datacenter replication** support

---

### Consistency Models

#### Strong Consistency
- **Linearizability** for all operations (default mode)
- **Serializable snapshot isolation** for transactions
- **Raft-based replication** ensuring majority agreement
- **Synchronous replication** to quorum of replicas
- **Leader-based reads** for strict consistency

#### Eventual Consistency
- **Tunable consistency levels** (ONE, QUORUM, ALL)
- **Asynchronous replication** for higher availability
- **Read repair** to fix inconsistencies on read path
- **Hinted handoff** for temporary node failures
- **Anti-entropy** with Merkle tree-based reconciliation
- **Vector clocks** for causality tracking
- **Conflict resolution strategies**:
  - Last-Write-Wins (LWW)
  - Application-defined merge functions
  - Multi-value returns for client resolution

---

### High Availability & Fault Tolerance

#### Failure Handling
- **Automatic failover** with sub-second detection
- **Split-brain prevention** via quorum requirements
- **Network partition tolerance** (CP in CAP theorem)
- **Graceful degradation** during partial failures
- **Automatic data repair** after node recovery
- **Backup and restore** utilities
- **Point-in-time recovery** (PITR)

#### Monitoring & Observability
- **Prometheus metrics** export
- **Real-time performance dashboards**
- **Detailed operation logging**
- **Distributed tracing** integration (OpenTelemetry)
- **Slow query logging** and analysis
- **Resource usage tracking** (CPU, memory, disk I/O)
- **Replication lag monitoring**
- **Cluster health status** API

---

### Performance Optimizations

#### Read Path
- **Multi-level caching**: Block cache, page cache, row cache
- **Bloom filters** to skip unnecessary SSTable reads
- **Index-only scans** when possible
- **Follower reads** for reduced leader load
- **Read-ahead** for sequential scans
- **Parallel query execution** for range scans

#### Write Path
- **Write buffering** in MemTable before disk flush
- **Batch writes** to reduce I/O operations
- **Group commit** for WAL efficiency
- **Background compaction** with rate limiting
- **Asynchronous replication** option
- **Zero-copy transfers** where possible

#### Query Optimization
- **Cost-based query planning**
- **Index selection** based on statistics
- **Predicate pushdown** to storage layer
- **Iterator pipelining** for memory efficiency
- **Parallel execution** for large scans

---

## Architecture

### System Components

```
┌─────────────────────────────────────────────────────────────┐
│                        Client Layer                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Go Client  │  │  HTTP/REST   │  │     gRPC     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Coordinator Layer                       │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Request Router  │  Query Planner  │  TX Manager   │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                    Distributed Layer                         │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │ Raft Leader  │  │   Follower   │  │   Follower   │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
│  ┌────────────────────────────────────────────────────┐     │
│  │  Membership  │  Replication  │  Failure Detection │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
                              │
┌─────────────────────────────────────────────────────────────┐
│                      Storage Layer                           │
│  ┌──────────────┐              ┌──────────────┐             │
│  │   B-Tree     │              │   LSM-Tree   │             │
│  │   Engine     │              │   Engine     │             │
│  └──────────────┘              └──────────────┘             │
│  ┌────────────────────────────────────────────────────┐     │
│  │ WAL  │  Buffer Pool  │  Page Cache  │  Compaction │     │
│  └────────────────────────────────────────────────────┘     │
│  ┌────────────────────────────────────────────────────┐     │
│  │              File System / Disk I/O                 │     │
│  └────────────────────────────────────────────────────┘     │
└─────────────────────────────────────────────────────────────┘
```

## Development

### Building from Source

```bash
# Prerequisites
# - Go 1.21 or higher
# - Make
# - Protocol Buffers compiler

# Clone and build
git clone https://github.com/yourusername/anchor.git
cd anchor
make deps
make build

# Run tests
make test

# Run benchmarks
make bench

```

### Project Structure

```
anchor/
├── cmd/
│   ├── anchor/            # Server binary
│   ├── anchor-cli/        # CLI client
│   └── anchor-admin/      # Admin utilities
├── pkg/
│   ├── storage/           # Storage engines
│   │   ├── btree/         # B-Tree implementation
│   │   ├── lsm/           # LSM-Tree implementation
│   │   └── page/          # Page format and buffer pool
│   ├── wal/               # Write-ahead log
│   ├── transaction/       # Transaction manager
│   ├── replication/       # Raft implementation
│   ├── cluster/           # Cluster membership
│   ├── consensus/         # Consensus protocols
│   └── client/            # Client library
├── internal/
│   ├── api/               # HTTP and gRPC APIs
│   ├── config/            # Configuration management
│   └── monitoring/        # Metrics and monitoring
├── tests/
│   ├── integration/       # Integration tests
│   ├── chaos/             # Chaos engineering tests
│   └── perf/              # Performance benchmarks
└── docs/                  # Documentation
```

---

## Roadmap

### Phase 1: Storage Foundation (In Progress)
- Basic key-value store
- Persistence layer
- Recovery mechanisms
- B-Tree implementation

### Phase 2: Advanced Storage
- LSM-Tree implementation
- Storage engine optimizations
- Page management
- Caching strategies

### Phase 3: Transactions
- MVCC and isolation
- WAL implementation
- Recovery algorithms
- Distributed transactions

### Phase 4: Distribution
- Cluster formation
- Failure detection
- Leader election
- Data partitioning

### Phase 5: Consistency & Replication
- Raft consensus
- Log replication
- Anti-entropy
- Conflict resolution

### Phase 6: Production Readiness
- Monitoring and metrics
- Operational tooling
- Performance optimization
- Security hardening
- Comprehensive documentation

---

## Implementation Checklist

Track the development progress of Anchor features:

### Storage Layer
- [x] Basic key-value operations (put, get, delete)
- [x] In-memory storage
- [x] Append-only log for persistence
- [ ] Recovery on restart
- [ ] B-Tree implementation
  - [ ] Node structure (internal vs leaf nodes)
  - [ ] Insertion with automatic splits
  - [ ] Search operations
  - [ ] Deletion operations
  - [ ] Disk-based serialization
  - [ ] Page cache (LRU)
- [ ] LSM-Tree implementation
  - [ ] MemTable (skiplist)
  - [ ] SSTable generation
  - [ ] Bloom filters
  - [ ] Compaction (leveled/tiered)
  - [ ] Block-based indexing
- [ ] Page format and management
  - [ ] Page header design
  - [ ] Slotted page structure
  - [ ] Checksum validation
  - [ ] Page compaction
  - [ ] Binary serialization

### Transaction Processing
- [ ] Write-Ahead Log (WAL)
  - [ ] LSN-based ordering
  - [ ] WAL writer and reader
  - [ ] Force-on-commit
  - [ ] Group commit optimization
- [ ] Transaction manager
  - [ ] Begin/commit/abort operations
  - [ ] MVCC implementation
  - [ ] Snapshot isolation
  - [ ] Lock management
- [ ] Recovery system
  - [ ] ARIES-style recovery
  - [ ] Analysis phase
  - [ ] Redo phase
  - [ ] Undo phase
  - [ ] Checkpointing

### B-Tree Variants & Optimizations
- [ ] Copy-on-Write B-Tree
  - [ ] Path copying
  - [ ] Version management
  - [ ] Garbage collection
  - [ ] Snapshot isolation
- [ ] B-Tree optimizations
  - [ ] Sibling pointers
  - [ ] Bulk loading
  - [ ] Prefix compression
  - [ ] Concurrent access (locks)

### Distributed System
- [ ] Cluster management
  - [ ] Node discovery
  - [ ] Static membership
  - [ ] Heartbeat mechanism
  - [ ] Configuration service
- [ ] Failure detection
  - [ ] Phi-accrual detector
  - [ ] Gossip protocol (SWIM)
  - [ ] Suspicion mechanism
  - [ ] Indirect probing
- [ ] Leader election
  - [ ] Raft election algorithm
  - [ ] Term numbers
  - [ ] Vote request/response
  - [ ] Randomized timeouts
  - [ ] Pre-vote optimization

### Replication & Consistency
- [ ] Raft consensus
  - [ ] Log replication
  - [ ] Append entries RPC
  - [ ] Log matching
  - [ ] Commit index advancement
  - [ ] Log compaction (snapshots)
- [ ] Strong consistency
  - [ ] Linearizable reads
  - [ ] Quorum-based operations
  - [ ] Leader-based routing
- [ ] Eventual consistency
  - [ ] Anti-entropy
  - [ ] Merkle trees
  - [ ] Read repair
  - [ ] Hinted handoff
  - [ ] Vector clocks
  - [ ] Conflict resolution

### Distributed Transactions
- [ ] Two-phase commit (2PC)
  - [ ] Transaction coordinator
  - [ ] Prepare phase
  - [ ] Commit phase
  - [ ] Timeout handling
  - [ ] Coordinator recovery
- [ ] Distributed deadlock handling
  - [ ] Deadlock detection
  - [ ] Timeout-based prevention
- [ ] Cross-partition transactions
  - [ ] Snapshot isolation
  - [ ] Optimistic concurrency control

### Data Distribution
- [ ] Partitioning
  - [ ] Hash-based sharding
  - [ ] Range-based partitioning
  - [ ] Consistent hashing
  - [ ] Virtual nodes
- [ ] Rebalancing
  - [ ] Automatic rebalancing
  - [ ] Data migration
  - [ ] Load-aware placement

### Client & API
- [ ] Client library
  - [ ] Connection management
  - [ ] Request routing
  - [ ] Retry logic
  - [ ] Connection pooling
- [ ] REST API
  - [ ] Basic CRUD operations
  - [ ] Transaction endpoints
  - [ ] Range queries
  - [ ] Admin operations
- [ ] gRPC API
  - [ ] Protocol buffer definitions
  - [ ] Service implementation
  - [ ] Streaming support

### Operations & Monitoring
- [ ] CLI tools
  - [ ] Database CLI client
  - [ ] Admin utilities
  - [ ] Debug tools
- [ ] Monitoring
  - [ ] Metrics collection
  - [ ] Prometheus integration
  - [ ] Dashboard
  - [ ] Health checks
- [ ] Backup & Restore
  - [ ] Snapshot creation
  - [ ] Backup utilities
  - [ ] Point-in-time recovery
- [ ] Cluster operations
  - [ ] Add/remove nodes
  - [ ] Membership changes
  - [ ] Rolling upgrades

### Testing & Quality
- [ ] Unit tests
  - [ ] Storage layer tests
  - [ ] Transaction tests
  - [ ] Replication tests
- [ ] Integration tests
  - [ ] End-to-end scenarios
  - [ ] Multi-node tests
- [ ] Chaos testing
  - [ ] Failure injection
  - [ ] Network partitions
  - [ ] Byzantine faults
- [ ] Performance benchmarks
  - [ ] Throughput tests
  - [ ] Latency measurements
  - [ ] Scalability tests

### Documentation
- [ ] API documentation
- [ ] Architecture guide
- [ ] Operations manual
- [ ] Performance tuning guide
- [ ] Troubleshooting guide

---