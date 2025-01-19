FROM golang:1.22.5

WORKDIR /planner-golang

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . . 

WORKDIR /planner-golang/cmd/planner-golang

RUN go build -o /planner-golang/bin/planner-golang .

EXPOSE 8080
ENTRYPOINT [ "/planner-golang/bin/planner-golang" ]