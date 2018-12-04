FROM scratch
ADD build/linux/processengine /
CMD ["/processengine", "-s", "nats://nats.nats.svc.cluster.local:4222", "-p", "/config/process.json"]