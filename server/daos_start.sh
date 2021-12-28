#/bin/bash

# start daos_server
sudo /opt/daos/bin/daos_server start -o /home/daos/daos_server.yml &

# wait a bit for the server starting
sleep 5
# format PMEM
dmg -i storage format

wait
