#!/bin/bash

echo $FASTLY_EXPORTER_SERVICES | tr "," "\n" > tmpfile

for i in $(cat tmpfile); do SERVICES="$SERVICES -service $i"; done

args="-token ${FASTLY_API_TOKEN} -endpoint http://0.0.0.0:${FASTLY_EXPORTER_PORT}/metrics $SERVICES"

if [ ${FASTLY_NAMESPACE} != "" ]; then
    args="$args -namespace ${FASTLY_NAMESPACE}"
fi

fastly-exporter $args


