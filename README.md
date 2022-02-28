### Golang CQRS EventSourcing EventStoreDB MongoDB Elasticsearch gRPC microservice example 👋

#### 👨‍💻 Full list what has been used:
[EventStoreDB](https://www.eventstore.com/) The database built for Event Sourcing<br/>
[gRPC](https://github.com/grpc/grpc-go) Go implementation of gRPC<br/>
[Jaeger](https://www.jaegertracing.io/) open source, end-to-end distributed [tracing](https://opentracing.io/)<br/>
[Prometheus](https://prometheus.io/) monitoring and alerting<br/>
[Grafana](https://grafana.com/) for to compose observability dashboards with everything from Prometheus<br/>
[MongoDB](https://github.com/mongodb/mongo-go-driver) Web and API based SMTP testing<br/>
[Elasticsearch](https://github.com/olivere/elastic) Elasticsearch client for Go.<br/>
[Redis](https://github.com/go-redis/redis) Type-safe Redis client for Golang<br/>
[swag](https://github.com/swaggo/swag) Swagger for Go<br/>
[Echo](https://github.com/labstack/echo) web framework<br/>
[Kibana](https://github.com/labstack/echo) Kibana is user interface that lets you visualize your Elasticsearch<br/>

### EventStoreDB UI:

http://localhost:2113

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3005

### Swagger UI:

http://localhost:5007/swagger/index.html

### Kibana UI:

http://localhost:5601/app/home#/


For local development 🙌👨‍💻🚀:

```
make local // for run docker compose
make run_es // run microservice
```
or 
```
make dev // run all in docker compose with hot reload
```