FROM golang:alpine
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
WORKDIR /build
#RUN git clone git@github.com:YasiruR/ktool-backend.git
COPY . .
#RUN ls -a
RUN go clean && go build -o ktool-backend

#WORKDIR /app
#RUN ls /build/ -al
#RUN cp /build/ktool-backend .
EXPOSE 7070
ENTRYPOINT ./ktool-backend
