FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN cd backend/server && go build -o /app/aether

FROM alpine:3.18
COPY --from=builder /app/aether /bin/aether
# copy frontend dist if present
COPY frontend/dist /var/www/dist
EXPOSE 8080
CMD ["/bin/aether"]