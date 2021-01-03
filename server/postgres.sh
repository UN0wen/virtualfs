#! /bin/bash

psql postgresql://postgres:postgres@localhost:5432/virtualfs $1 $2
