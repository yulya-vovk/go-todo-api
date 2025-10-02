#!/bin/bash
curl -X PUT http://localhost:8080/tasks/2 \
  -H "Content-Type: application/json" \
  -d '{"title":"Задача обновлена! ","done":true}'