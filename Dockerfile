FROM scratch
ADD docker112rc3-runtimefix /bin/docker112rc3-runtimefix
ENTRYPOINT ["/bin/docker112rc3-runtimefix"]
CMD ["/docker"]
