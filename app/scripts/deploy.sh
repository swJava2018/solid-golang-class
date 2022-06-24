#!/usr/bin/env bash

overrideValues="${1}"
override=""
if [[ -n "${overrideValues}" ]]; then
    override="-f ${overrideValues}"
fi
service="event-data-pipeline"
namespace=${NAMESPACE=default}
releaseName=${RELEASENAME=$service}
helm upgrade --install --atomic $releaseName ../helm -f ../helm/values.yaml "$override" -n $namespace