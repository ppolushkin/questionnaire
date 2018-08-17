FROM iron/go

WORKDIR /run

ADD ./questionnaire /run

EXPOSE 8080

CMD ["/run/questionnaire"]

#ENTRYPOINT ["./questionnaire"]