ARG DEV_IMG=dev_vpp_agent
FROM ${DEV_IMG} as devimg

FROM ubuntu:18.04 as base

RUN apt-get update && apt-get install -y --no-install-recommends \
		# general tools
		inetutils-traceroute \
		iproute2 \
		iputils-ping \
		# vpp requirements
		ca-certificates \
		libapr1 \
		libc6 \
		libmbedcrypto1 \
		libmbedtls10 \
		libmbedx509-0 \
		libnuma1 \
		openssl \
 	&& rm -rf /var/lib/apt/lists/*

# install vpp
COPY --from=devimg /vpp/*.deb /opt/vpp/

RUN set -eux; \
	cd /opt/vpp/; \
	apt-get update; \
	apt-get install -y ./*.deb; \
	rm *.deb; \
	rm -rf /var/lib/apt/lists/*;

# Copy configs
COPY \
	etcd.conf \
	grpc.conf \
	supervisor.conf \
 /opt/vpp-agent/dev/

COPY vpp.conf /etc/vpp/vpp.conf
COPY init_hook.sh /usr/bin/

# handle differences in vpp.conf which are between supported VPP versions
ARG VPP_VERSION
COPY legacy-nat.conf /tmp/legacy-nat.conf
RUN bash -c "if [[ \"$VPP_VERSION\" < "2101" ]]; then cat /tmp/legacy-nat.conf >> /etc/vpp/vpp.conf; fi"
RUN rm /tmp/legacy-nat.conf

# Install agent
COPY --from=devimg \
    /go/bin/agentctl \
    /go/bin/vpp-agent \
    /go/bin/vpp-agent-init \
 /bin/

# Final image
FROM scratch
COPY --from=base / /

WORKDIR /root/

ENV SUPERVISOR_CONFIG=/opt/vpp-agent/dev/supervisor.conf

CMD rm -f /dev/shm/db /dev/shm/global_vm /dev/shm/vpe-api && \
	mkdir -p /run/vpp && \
	exec vpp-agent-init
