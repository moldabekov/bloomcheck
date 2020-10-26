FROM scratch

EXPOSE 9998

COPY bloomcheck /

COPY bloom.filter /

CMD ["/bloomcheck", "/bloom.filter"]
