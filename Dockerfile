FROM amazonlinux:2

RUN yum install -y zip && yum clean all

WORKDIR /app

ARG FUNCTION_NAME
COPY --from=builder /app/bootstrap /app/bootstrap
RUN zip -j /app/${FUNCTION_NAME}_function.zip /app/bootstrap

RUN mkdir -p /app/output
RUN cp /app/${FUNCTION_NAME}_function.zip /app/output/${FUNCTION_NAME}_function.zip
