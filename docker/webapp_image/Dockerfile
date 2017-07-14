FROM debian
LABEL version=0.4

ENV WA_DOCUMENT_DIR /webapp/docs
ENV WA_WEB_PORT 33000
ENV CRUD_HOST crud
ENV CRUD_PORT 33002
ENV SP_HOST stockprovidersvc.default
ENV SP_PORT 33001
ENV PATH /webapp:${PATH}

EXPOSE 33000

RUN ["mkdir","/webapp"]
RUN ["mkdir","/webapp/docs"]

COPY webapp /webapp
COPY docs /webapp/docs


WORKDIR /webapp

CMD ["webapp"]

