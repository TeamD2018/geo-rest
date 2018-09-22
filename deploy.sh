echo "$DOCKERHUB_PASSWORD" | docker login -u "$DOCKERHUB_NAME" --password-stdin
docker tag swagger-ui-geo-rest $DOCKERHUB_NAME/swagger-ui-geo-rest
docker push $DOCKERHUB_NAME/swagger-ui-geo-rest