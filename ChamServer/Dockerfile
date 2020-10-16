FROM python:3.6
RUN apt-get update
RUN apt-get install -y m4
RUN apt-get install -y flex
RUN apt-get install -y bison
RUN mkdir /var/Working
ADD ./gmp-6.2.0/ /var/Working/gmp
WORKDIR /var/Working/gmp
RUN chmod +x configure
RUN ./configure
RUN make
RUN make install
ADD ./pbc-0.5.14/ /var/Working/pbc
WORKDIR /var/Working/pbc
RUN chmod +x configure
RUN ./configure
RUN make
RUN make install
ADD ./openssl-1.1.1g/ /var/Working/openssl
WORKDIR /var/Working/openssl
RUN chmod +x config
RUN ./config
RUN make
RUN make install
RUN apt-get install -y python3-setuptools
RUN apt-get install -y python3-dev
RUN apt-get install -y libssl-dev
ADD ./charm-dev/ /var/Working/charm
WORKDIR /var/Working/charm
RUN chmod +x configure.sh
RUN ./configure.sh
RUN make
RUN make install
COPY . /home/python/mapchdocker
WORKDIR /home/python/mapchdocker
ENTRYPOINT [ "python" ]
CMD [ "MAPCH.py" ]