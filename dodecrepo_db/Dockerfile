FROM postgres

COPY setup.sql /tmp/setup.sql
COPY setup.sh /docker-entrypoint-initdb.d/
COPY docker-entrypoint-wrapper.sh /

RUN chmod +x /docker-entrypoint-wrapper.sh

ENTRYPOINT ["/docker-entrypoint-wrapper.sh"]
CMD ["postgres"]
