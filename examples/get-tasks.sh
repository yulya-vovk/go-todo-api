#!/bin/bash
curl -X GET http://localhost:8080/tasks \
  -H "Content-Type: application/json"