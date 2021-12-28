# Before build this container, build daos container image
# according to https://docs.daos.io/admin/installation/#building-a-docker-image
FROM daos:latest

RUN sudo yum update
RUN sudo yum install -y wget fuse3
RUN sudo yum install -y vim
RUN sudo yum remove -y golang
RUN cd /home/daos && wget https://go.dev/dl/go1.17.5.linux-amd64.tar.gz && tar xf go1.17.5.linux-amd64.tar.gz
RUN echo "export PATH=$PATH:/home/daos/go/bin" >> /home/daos/.bashrc
RUN source /home/daos/.bashrc && cd /home/daos && git clone https://github.com/rexray/gocsi && cd gocsi/ && make

COPY ["daos_agent.yml", "/home/daos/"]
COPY ["daos_control.yml", "/opt/daos/etc/daos_control.yml"]
COPY ["daos_start.sh", "/home/daos/"]
COPY ["cmd/dfuseplugin/main", "/dfuseplugin"]
ENTRYPOINT ["/bin/bash", "/home/daos/daos_start.sh"]

