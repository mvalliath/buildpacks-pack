FROM golang:windowsservercore-1809 AS builder

WORKDIR /workdir
COPY . .
ENV GO111MODULE on
RUN go build -mod=vendor -o c:/workdir/phase.exe ./phase.go

FROM mcr.microsoft.com/windows/nanoserver:1809

RUN mkdir c:\lifecycle
RUN mkdir c:\buildpacks
COPY --from=builder /workdir/phase.exe /lifecycle

RUN echo "original-order-toml" > c:\buildpacks\order.toml

ENV CNB_USER_ID 111
ENV CNB_GROUP_ID 222

LABEL io.buildpacks.stack.id=com.example.stack