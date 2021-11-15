FROM golang as builder

WORKDIR /go/src/usermanagement
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 go build -v

FROM alpine
LABEL maintainer="jayaraj.esvar@gmail.com"
WORKDIR /home
COPY --from=builder /go/src/usermanagement/usermanagement /home
COPY --from=builder /go/src/usermanagement/.usermanagement.yml /home/.usermanagement.yml
CMD [ "/home/usermanagement" ]
