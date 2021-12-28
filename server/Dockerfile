FROM daos:latest

RUN sudo yum install -y wget fuse3
RUN wget -P /opt/daos/examples/ https://raw.githubusercontent.com/daos-stack/daos/release/1.2/utils/config/daos_agent.yml

COPY ["daos_server.yml", "/home/daos/"]
COPY ["daos_start.sh", "/home/daos/"]
ENTRYPOINT ["/bin/bash", "/home/daos/daos_start.sh"]

