FROM golang:1.10
EXPOSE 8080
ADD tinyworker /
WORKDIR /
CMD ["/tinyworker"]