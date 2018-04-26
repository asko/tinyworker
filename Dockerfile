FROM scratch
EXPOSE 8080
ADD tinyworker /
WORKDIR /
CMD ["/tinyworker"]