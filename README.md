### Golang CQRS EventSourcing EventStoreDB MongoDB gRPC microservice example 👋

#### 👨‍💻 Full list what has been used:
[EventStoreDB](https://www.eventstore.com/) The database built for Event Sourcing<br/>
[gRPC](https://github.com/grpc/grpc-go) Go implementation of gRPC<br/>
[Jaeger](https://www.jaegertracing.io/) open source, end-to-end distributed [tracing](https://opentracing.io/)<br/>
[Prometheus](https://prometheus.io/) monitoring and alerting<br/>
[Grafana](https://grafana.com/) for to compose observability dashboards with everything from Prometheus<br/>
[MongoDB](https://github.com/mongodb/mongo-go-driver) Web and API based SMTP testing<br/>
[Redis](https://github.com/go-redis/redis) Type-safe Redis client for Golang<br/>
[swag](https://github.com/swaggo/swag) Swagger for Go<br/>
[Echo](https://github.com/labstack/echo) web framework<br/>

### EventStoreDB UI:

http://localhost:2113

### Jaeger UI:

http://localhost:16686

### Prometheus UI:

http://localhost:9090

### Grafana UI:

http://localhost:3000

### Swagger UI:

http://localhost:5007/swagger/index.html


For local development 🙌👨‍💻🚀:

```
make mongo // run mongo init scripts
make swagger // generate swagger documentation
make local or docker_dev // for run docker compose files
```