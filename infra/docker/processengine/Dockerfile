FROM scratch
ADD data/test/process.json /
ADD build/linux/processengine /
CMD ["/processengine", "-s", "nats://nats.nats.svc.cluster.local:4222", "-p", "./process.json"]