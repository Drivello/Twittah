#!/bin/bash

set -e

KAFKA_BROKER=${KAFKA_BROKER:-localhost:9092}
REPLICATION_FACTOR=${REPLICATION_FACTOR:-1}
PARTITIONS=${PARTITIONS:-1}

create_topic() {
  local topic=$1
  /opt/bitnami/kafka/bin/kafka-topics.sh \
    --bootstrap-server "$KAFKA_BROKER" \
    --create \
    --if-not-exists \
    --topic "$topic" \
    --replication-factor "$REPLICATION_FACTOR" \
    --partitions "$PARTITIONS"
}

create_topic users.events
create_topic follows.events
create_topic tweets.events
create_topic timelines.updated

echo "Kafka topics initialized."
