#!/bin/bash

docker exec -it -u postgres postgres_database bash -c "cd /var/lib/postgresql/data && chmod 777 ."
