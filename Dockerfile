FROM scratch
EXPOSE 8080
ENTRYPOINT ["/blocktrackmeapi"]
COPY ./bin/ /