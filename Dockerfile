FROM scratch
EXPOSE 8080
ENTRYPOINT ["/goqs"]
COPY ./bin/ /