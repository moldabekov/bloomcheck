FROM scratch

EXPOSE 9998

COPY bloomcheck /

COPY bloom.filter /

CMD ["/usr/local/bin/bloomcheck", "/bloom.filter"]
