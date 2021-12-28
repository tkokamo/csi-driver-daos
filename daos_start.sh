#/bin/bash

# start daos_agent
sudo /opt/daos/bin/daos_agent -i -o /home/daos/daos_agent.yml &
sudo /dfuseplugin "$@"
