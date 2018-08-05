FROM golang

# Install dep
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com/ramonhpr
WORKDIR /go/src/github.com/ramonhpr

RUN git clone https://github.com/ramonhpr/bleu_hackathon.git
WORKDIR /go/src/github.com/ramonhpr/bleu_hackathon

RUN go get -u github.com/aws/aws-sdk-go/...
RUN go get -d github.com/ethereum/go-ethereum

RUN go build -o dona-maria .

ENV AWS_SDK_LOAD_CONFIG=true
ENV AWS_REGION=us-east-2

RUN mkdir -p ~/.aws/

RUN echo "[bleu-hackathon]\naws_access_key_id = \naws_secret_access_key = " > ~/.aws/credentials

EXPOSE 8080

CMD ["./dona-maria"]

