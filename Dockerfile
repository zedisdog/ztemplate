FROM ip:port/library/gobuilder:v0.0.2 AS builder //自制gobuilder镜像

ARG GIT_USERNAME
ARG GIT_PASSWORD
ARG RUN_ENV

# http方式的仓库自动输账户密码的配置
# RUN echo "http://${GIT_USERNAME}:${GIT_PASSWORD}@domain" > /root/.git-credentials &&\
#     git config --global credential.http://domain.helper store

WORKDIR /build

COPY . .
COPY ./etc/${RUN_ENV}.yaml /app/etc/conf.yaml
RUN go build -ldflags="-s -w" -o /app/app .

FROM scratch

COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /usr/share/zoneinfo/Asia/Shanghai
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/app /app/app
COPY --from=builder /app/etc /app/etc

CMD ["./app", "-f", "etc/conf.yaml"]


