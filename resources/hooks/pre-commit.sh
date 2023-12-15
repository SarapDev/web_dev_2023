#!/bin/bash

# запускаем проверку как на github workflow
# (см. .github/workflows/test.yml)
make check || exit 1

# прогоняем весь проект через go fmt
go fmt ./... || exit 1
