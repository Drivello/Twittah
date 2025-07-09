#!/bin/bash

set -e

KAFKA_BROKER=${KAFKA_BROKER:-localhost:9092}
REPLICATION_FACTOR=${REPLICATION_FACTOR:-1}
PARTITIONS=${PARTITIONS:-1}

create_topic() {
  local topic=$1
  kafka-topics \
    --bootstrap-server "$KAFKA_BROKER" \
    --create \
    --if-not-exists \
    --topic "$topic" \
    --replication-factor "$REPLICATION_FACTOR" \
    --partitions "$PARTITIONS"
}

create_topic tweets.published
create_topic follows.created
create_topic follows.deleted
create_topic timelines.updated

echo "Kafka topics initialized."
